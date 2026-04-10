package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReminderStatus string

const (
	ReminderStatusScheduled ReminderStatus = "scheduled"
)

type Reminder struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TodoID    primitive.ObjectID `json:"todo_id" bson:"todo_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	RemindAt  time.Time          `json:"remind_at" bson:"remind_at"`
	Status    ReminderStatus     `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
