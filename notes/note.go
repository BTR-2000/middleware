package notes

import (
	"fmt"
	"sync"
	"time"
)

type Note struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Time        time.Time `json:"time"`
}

type NoteStore struct {
	mu     sync.RWMutex
	notes  map[int]Note
	nextID int
}

func NewNoteStore() *NoteStore {
	return &NoteStore{
		mu:     sync.RWMutex{},
		notes:  make(map[int]Note),
		nextID: 1,
	}
}

func (n *NoteStore) AllNotes() []Note {
	n.mu.RLock()
	defer n.mu.RUnlock()
	slice := make([]Note, 0, len(n.notes))

	for _, note := range n.notes {
		slice = append(slice, note)
	}
	return slice
}

func (n *NoteStore) CreateNote(newTitle, newDescription string) (Note, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if newTitle == "" || newDescription == "" {
		return Note{}, fmt.Errorf("Данные заполнены неверно!")
	}

	note := Note{
		Id:          n.nextID,
		Title:       newTitle,
		Description: newDescription,
		Time:        time.Now(),
	}

	n.notes[n.nextID] = note
	n.nextID++
	return note, nil
}

func (n *NoteStore) DeleteByID(id int) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, exist := n.notes[id]; !exist {
		return fmt.Errorf("Пользователя с id = %d не существует!", id)
	}

	delete(n.notes, id)
	return nil
}
