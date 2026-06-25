package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"newMiddle/notes"
	"strconv"
)

type Server struct {
	store *notes.NoteStore
}

func NewServer(newStore *notes.NoteStore) *Server {
	return &Server{
		store: newStore,
	}
}

func (s *Server) GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	notes := s.store.AllNotes()
	if len(notes) == 0 {
		w.Write([]byte("Пока заметок нет"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "Ошибка передачи заметок", http.StatusInternalServerError)
		return
	}
}

func (s *Server) NewNoteHandler(w http.ResponseWriter, r *http.Request) {
	var noteDTO notes.NoteDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&noteDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "неверные данные json"})
		return
	}

	newNote, err := s.store.CreateNote(noteDTO.Title, noteDTO.Description)
	if err != nil {
		strError := fmt.Sprintf("%s", err)
		http.Error(w, strError, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newNote)
}

func (s *Server) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	n := r.PathValue("id")
	id, err := strconv.Atoi(n)
	if err != nil {
		http.Error(w, "Ошибка id", http.StatusBadRequest)
		return
	}

	if err = s.store.DeleteByID(id); err != nil {
		strError := fmt.Sprintf("%s", err)
		http.Error(w, strError, http.StatusNotFound)
		return
	}
	w.Write([]byte("Данные успешно удалены!"))
}
