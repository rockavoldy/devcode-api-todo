package model

import (
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrTitleNull           = errors.New("title cannot be null")
	ErrActivityGroupIdNull = errors.New("activity_group_id cannot be null")
	ErrRecordNotFound      = errors.New("record not found")
)

var ActivityIncr int = 0

func GetActivityId(m *sync.Mutex) int {
	m.Lock()

	ActivityIncr += 1
	m.Unlock()

	return ActivityIncr
}

type ActivityGroup struct {
	ID        int       `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Title     string    `json:"title,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
}

type PrintActivtyGroup struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type PrintActivityGroups struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

func NewActivityGroup(m *sync.Mutex, email, title string) *ActivityGroup {
	return &ActivityGroup{
		ID:        GetActivityId(m),
		Email:     email,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (a *ActivityGroup) MapToInterface() map[string]interface{} {
	if a == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"id":         a.ID,
		"email":      a.Email,
		"title":      a.Title,
		"created_at": a.CreatedAt,
		"updated_at": a.UpdatedAt,
		"deleted_at": a.DeletedAt,
	}
}

func (a *ActivityGroup) Validate() error {
	if a.Title == "" {
		return ErrTitleNull
	}

	return nil
}

type Priority int

const (
	VeryHigh Priority = iota
	High
	Normal
	Low
	VeryLow
)

func (p Priority) String() string {
	priorityString := map[Priority]string{
		0: "very-high",
		1: "high",
		2: "normal",
		3: "low",
		4: "very-low",
	}

	return priorityString[p]
}

func PriorityStringToInt(priority string) Priority {
	priorityInt := map[string]Priority{
		"very-high": 0,
		"high":      1,
		"normal":    2,
		"low":       3,
		"very-low":  4,
	}

	return priorityInt[priority]
}

var TodoIncr int = 0

func GetTodoId(m *sync.Mutex) int {
	m.Lock()

	TodoIncr += 1
	m.Unlock()

	return TodoIncr
}

type TodoItem struct {
	ID              int       `json:"id,omitempty"`
	ActivityGroupId int       `json:"activity_group_id,omitempty"`
	Title           string    `json:"title,omitempty"`
	IsActive        bool      `json:"is_active,omitempty"`
	Priority        string    `json:"priority,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	DeletedAt       time.Time `json:"deleted_at,omitempty"`
}

type PrintTodoItem struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type PrintTodoItems struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

func NewTodoItem(m *sync.Mutex, activity_group_id int, title string, isActive bool, priority Priority) *TodoItem {
	if !isActive {
		isActive = true
	}
	return &TodoItem{
		ID:              GetTodoId(m),
		ActivityGroupId: activity_group_id,
		Title:           title,
		IsActive:        isActive,
		Priority:        priority.String(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func (t *TodoItem) MapToInterface() map[string]interface{} {
	if t == nil {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"id":                t.ID,
		"activity_group_id": t.ActivityGroupId,
		"title":             t.Title,
		"is_active":         t.IsActive,
		"priority":          t.Priority,
		"created_at":        t.CreatedAt,
		"updated_at":        t.UpdatedAt,
		"deleted_at":        t.DeletedAt,
	}
}

func (t *TodoItem) Validate() error {
	if t.Title == "" {
		return ErrTitleNull
	}

	if t.ActivityGroupId == 0 {
		return ErrActivityGroupIdNull
	}

	return nil
}
