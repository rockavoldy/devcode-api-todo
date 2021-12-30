package main

import (
	"github.com/pkg/errors"
)

var (
	ErrTitleNull           = errors.New("title cannot be null")
	ErrActivityGroupIdNull = errors.New("activity_group_id cannot be null")
	ErrRecordNotFound      = errors.New("record not found")
)

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
