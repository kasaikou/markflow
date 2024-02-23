package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNotCirculated(t *testing.T) {
	document := Document{
		Tasks: map[string]DocumentTask{
			"task-1-1": {
				DependTasks: []string{"task-2"},
			},
			"task-1-2": {
				DependTasks: []string{"task-3-1", "task-4"},
			},
			"task-2": {
				DependTasks: []string{"task-3-1", "task-3-2"},
			},
			"task-3-1": {
				DependTasks: []string{"task-4"},
			},
			"task-3-2": {
				DependTasks: []string{"task-4"},
			},
			"task-4": {},
		},
	}

	err := validateIsTaskDependencyCirculated(document)
	assert.NoError(t, err)
}

func TestValidateCirculated(t *testing.T) {
	document := Document{
		Tasks: map[string]DocumentTask{
			"task-1-1": {
				DependTasks: []string{"task-2"},
			},
			"task-1-2": {
				DependTasks: []string{"task-3-1", "task-4"},
			},
			"task-2": {
				DependTasks: []string{"task-3-1", "task-3-2"},
			},
			"task-3-1": {
				DependTasks: []string{"task-4"},
			},
			"task-3-2": {
				DependTasks: []string{"task-4"},
			},
			"task-4": {
				DependTasks: []string{"task-2"},
			},
		},
	}

	err := validateIsTaskDependencyCirculated(document)
	assert.Error(t, err)
}
