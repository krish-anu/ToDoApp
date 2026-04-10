package workflow

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/krish-anu/ToDoAppBackend/internal/ai"
	"github.com/krish-anu/ToDoAppBackend/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidMessage  = errors.New("message is required")
	ErrInvalidTimezone = errors.New("timezone is invalid")
)

type TaskFromTextInput struct {
	UserID   string
	Message  string
	Timezone string
}

type TaskFromTextResult struct {
	Todo     models.Todo       `json:"todo"`
	Reminder *models.Reminder  `json:"reminder,omitempty"`
	Parsed   ai.TaskExtraction `json:"parsed"`
	Partial  bool              `json:"partial"`
	Message  string            `json:"message"`
}

type TaskFromTextWorkflow struct {
	todoCollection         *mongo.Collection
	reminderCollection     *mongo.Collection
	workflowRunsCollection *mongo.Collection
	extractor              ai.TaskExtractor
	defaultUserID          string
	defaultTimezone        string
}

func NewTaskFromTextWorkflow(
	todoCollection *mongo.Collection,
	reminderCollection *mongo.Collection,
	workflowRunsCollection *mongo.Collection,
	extractor ai.TaskExtractor,
	defaultUserID string,
	defaultTimezone string,
) *TaskFromTextWorkflow {
	resolvedUserID := strings.TrimSpace(defaultUserID)
	if resolvedUserID == "" {
		resolvedUserID = "local-user"
	}

	resolvedTimezone := strings.TrimSpace(defaultTimezone)
	if resolvedTimezone == "" {
		resolvedTimezone = "UTC"
	}

	return &TaskFromTextWorkflow{
		todoCollection:         todoCollection,
		reminderCollection:     reminderCollection,
		workflowRunsCollection: workflowRunsCollection,
		extractor:              extractor,
		defaultUserID:          resolvedUserID,
		defaultTimezone:        resolvedTimezone,
	}
}

func (w *TaskFromTextWorkflow) Run(ctx context.Context, input TaskFromTextInput) (result *TaskFromTextResult, err error) {
	message := strings.TrimSpace(input.Message)
	if message == "" {
		return nil, ErrInvalidMessage
	}

	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		userID = w.defaultUserID
	}

	timezone := strings.TrimSpace(input.Timezone)
	if timezone == "" {
		timezone = w.defaultTimezone
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidTimezone, timezone)
	}

	startedAt := time.Now().UTC()
	steps := make([]workflowStep, 0, 6)
	addStep := func(name string, status string, detail string) {
		steps = append(steps, workflowStep{
			Name:      name,
			Status:    status,
			Detail:    detail,
			Timestamp: time.Now().UTC(),
		})
	}

	var createdTodoID *primitive.ObjectID
	var createdReminderID *primitive.ObjectID
	var runError string

	defer func() {
		if err != nil {
			runError = err.Error()
		}

		runStatus := "success"
		if err != nil {
			runStatus = "failed"
		} else if result != nil && result.Partial {
			runStatus = "partial_success"
		}

		w.saveRun(ctx, workflowRunDocument{
			WorkflowName:  "task_from_text",
			UserID:        userID,
			InputMessage:  message,
			Timezone:      timezone,
			Status:        runStatus,
			Error:         runError,
			StartedAt:     startedAt,
			CompletedAt:   time.Now().UTC(),
			Steps:         steps,
			CreatedTodoID: createdTodoID,
			ReminderID:    createdReminderID,
			CreatedAt:     startedAt,
		})
	}()

	addStep("extract_task", "running", "Extracting task details from user message")
	extractedTask, err := w.extractor.ExtractTaskFromText(ctx, message, timezone)
	if err != nil {
		addStep("extract_task", "failed", err.Error())
		return nil, err
	}
	addStep("extract_task", "success", fmt.Sprintf("title=%q", extractedTask.Title))

	title := strings.TrimSpace(extractedTask.Title)
	if title == "" {
		addStep("validate_extraction", "failed", "title is empty")
		return nil, errors.New("extracted title is empty")
	}

	dueAt, err := parseOptionalDateTime(extractedTask.DueAt, location)
	if err != nil {
		addStep("validate_extraction", "failed", "invalid due_at")
		return nil, fmt.Errorf("invalid due_at value: %w", err)
	}

	reminderAt, reminderParseErr := parseOptionalDateTime(extractedTask.RemindAt, location)
	partial := false
	resultMessage := "Task created successfully."

	if reminderParseErr != nil {
		partial = true
		resultMessage = "Task created, but reminder was skipped because reminder time was invalid."
		addStep("validate_extraction", "warning", "invalid remind_at; reminder will be skipped")
	} else {
		addStep("validate_extraction", "success", "extracted data is valid")
	}

	now := time.Now().UTC()
	todo := models.Todo{
		Body:      title,
		Completed: false,
		Priority:  normalizePriority(extractedTask.Priority),
		Tags:      normalizeTags(extractedTask.Tags),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if dueAt != nil {
		todo.DueAt = dueAt
	}
	if reminderAt != nil {
		todo.RemindAt = reminderAt
	}

	addStep("create_task", "running", "Saving todo")
	insertResult, err := w.todoCollection.InsertOne(ctx, todo)
	if err != nil {
		addStep("create_task", "failed", err.Error())
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		addStep("create_task", "failed", "inserted id is not an ObjectID")
		return nil, errors.New("failed to parse inserted todo id")
	}

	todo.ID = insertedID
	createdTodoID = &todo.ID
	addStep("create_task", "success", fmt.Sprintf("todo_id=%s", todo.ID.Hex()))

	var reminder *models.Reminder
	if reminderAt != nil {
		switch {
		case w.reminderCollection == nil:
			partial = true
			resultMessage = "Task created, but reminders collection is not configured."
			addStep("schedule_reminder", "failed", "reminders collection is nil")
		case reminderAt.Before(time.Now().UTC()):
			partial = true
			resultMessage = "Task created, but reminder was skipped because reminder time is in the past."
			addStep("schedule_reminder", "skipped", "reminder time is in the past")
		default:
			addStep("schedule_reminder", "running", "Scheduling reminder")
			reminderRecord := models.Reminder{
				TodoID:    todo.ID,
				UserID:    userID,
				RemindAt:  reminderAt.UTC(),
				Status:    models.ReminderStatusScheduled,
				CreatedAt: time.Now().UTC(),
			}

			reminderInsertResult, reminderErr := w.reminderCollection.InsertOne(ctx, reminderRecord)
			if reminderErr != nil {
				partial = true
				resultMessage = "Task created, but reminder scheduling failed."
				addStep("schedule_reminder", "failed", reminderErr.Error())
			} else {
				if reminderID, castOK := reminderInsertResult.InsertedID.(primitive.ObjectID); castOK {
					reminderRecord.ID = reminderID
					createdReminderID = &reminderRecord.ID
				}
				reminder = &reminderRecord
				if !partial {
					resultMessage = "Task created and reminder scheduled."
				}
				addStep("schedule_reminder", "success", "reminder scheduled")
			}
		}
	} else {
		addStep("schedule_reminder", "skipped", "no reminder requested")
	}

	result = &TaskFromTextResult{
		Todo:     todo,
		Reminder: reminder,
		Parsed:   extractedTask,
		Partial:  partial,
		Message:  resultMessage,
	}
	return result, nil
}

func parseOptionalDateTime(raw *string, location *time.Location) (*time.Time, error) {
	if raw == nil {
		return nil, nil
	}

	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := parseDateTime(trimmed, location)
	if err != nil {
		return nil, err
	}

	utc := parsed.UTC()
	return &utc, nil
}

func parseDateTime(value string, location *time.Location) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, nil
	}

	layouts := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, location)
		if err != nil {
			continue
		}

		if layout == "2006-01-02" {
			parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 9, 0, 0, 0, location)
		}

		return parsed, nil
	}

	return time.Time{}, fmt.Errorf("unsupported datetime format: %s", value)
}

func normalizePriority(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "low", "medium", "high":
		return strings.ToLower(strings.TrimSpace(raw))
	default:
		return "medium"
	}
}

func normalizeTags(tags []string) []string {
	if len(tags) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(tags))
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleaned := strings.ToLower(strings.TrimSpace(tag))
		if cleaned == "" {
			continue
		}
		if _, exists := seen[cleaned]; exists {
			continue
		}
		seen[cleaned] = struct{}{}
		normalized = append(normalized, cleaned)
	}
	return normalized
}

func (w *TaskFromTextWorkflow) saveRun(ctx context.Context, run workflowRunDocument) {
	if w.workflowRunsCollection == nil {
		return
	}

	_, _ = w.workflowRunsCollection.InsertOne(ctx, run)
}

type workflowRunDocument struct {
	WorkflowName  string              `bson:"workflow_name"`
	UserID        string              `bson:"user_id"`
	InputMessage  string              `bson:"input_message"`
	Timezone      string              `bson:"timezone"`
	Status        string              `bson:"status"`
	Error         string              `bson:"error,omitempty"`
	StartedAt     time.Time           `bson:"started_at"`
	CompletedAt   time.Time           `bson:"completed_at"`
	Steps         []workflowStep      `bson:"steps"`
	CreatedTodoID *primitive.ObjectID `bson:"created_todo_id,omitempty"`
	ReminderID    *primitive.ObjectID `bson:"reminder_id,omitempty"`
	CreatedAt     time.Time           `bson:"created_at"`
}

type workflowStep struct {
	Name      string    `bson:"name"`
	Status    string    `bson:"status"`
	Detail    string    `bson:"detail"`
	Timestamp time.Time `bson:"timestamp"`
}
