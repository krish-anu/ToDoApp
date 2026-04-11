package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed" bson:"completed"`
	Body      string             `json:"body" bson:"body"`
	DueAt     *time.Time         `json:"due_at,omitempty" bson:"due_at,omitempty"`
	RemindAt  *time.Time         `json:"remind_at,omitempty" bson:"remind_at,omitempty"`
	Priority  string             `json:"priority,omitempty" bson:"priority,omitempty"`
	Tags      []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CreateTodoRequest struct {
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}

type UpdateTodoRequest struct {
	Completed *bool `json:"completed"`
}
