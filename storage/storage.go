package storage

import (
	"encoding/json"
	"lab1/workspace"
	"os"
)

// LocalStorage handles workspace state persistence
type LocalStorage struct {
	path string
}

// NewLocalStorage creates a new local storage instance
func NewLocalStorage(path string) *LocalStorage {
	return &LocalStorage{path: path}
}

// SaveMemento saves a workspace memento to disk
func (ls *LocalStorage) SaveMemento(memento *workspace.WorkspaceMemento) error {
	data, err := json.MarshalIndent(memento, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ls.path, data, 0644)
}

// LoadMemento loads a workspace memento from disk
func (ls *LocalStorage) LoadMemento() (*workspace.WorkspaceMemento, error) {
	data, err := os.ReadFile(ls.path)
	if err != nil {
		return nil, err
	}
	
	var memento workspace.WorkspaceMemento
	if err := json.Unmarshal(data, &memento); err != nil {
		return nil, err
	}
	
	return &memento, nil
}
