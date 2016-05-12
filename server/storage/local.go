package storage

import (
	"fmt"

	"github.com/glestaris/uberlist-server"
)

type LocalStore struct {
	tasks map[uberlist.TaskID]uberlist.Task
}

func NewLocalStore() *LocalStore {
	return &LocalStore{
		tasks: make(map[uberlist.TaskID]uberlist.Task),
	}
}

func (s *LocalStore) AddTask(task uberlist.Task) (uberlist.TaskID, error) {
	task.ID = uberlist.TaskID(len(s.tasks) + 1)
	s.tasks[task.ID] = task

	return task.ID, nil
}

func (s *LocalStore) UpdateTask(task uberlist.Task) error {
	_, ok := s.tasks[task.ID]
	if !ok {
		return fmt.Errorf("Task %d was not found!", task.ID)
	}

	s.tasks[task.ID] = task

	return nil
}

func (s *LocalStore) TaskByID(id uberlist.TaskID) (uberlist.Task, error) {
	t, ok := s.tasks[id]
	if !ok {
		return uberlist.Task{}, fmt.Errorf("Task %d was not found!", id)
	}

	return t, nil
}

func (s *LocalStore) OrderedTasks() ([]uberlist.Task, error) {
	retVal := make([]uberlist.Task, len(s.tasks))
	i := 0
	for _, t := range s.tasks {
		retVal[i] = t
		i++
	}

	return retVal, nil
}
