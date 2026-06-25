package main

import (
	"log"
	"net/http"
	"newMiddle/handlers"
	"newMiddle/middlewares"
	"newMiddle/notes"
)

func main() {
	store := notes.NewNoteStore()
	server := handlers.NewServer(store)

	mux := http.NewServeMux()

	mustHaveMiddlewares := middlewares.MustMiddlewares

	mux.Handle("GET /notes", mustHaveMiddlewares(http.HandlerFunc(server.GetAllNotesHandler)))
	mux.Handle("POST /notes", middlewares.ContentMiddleware(mustHaveMiddlewares(http.HandlerFunc(server.NewNoteHandler))))
	mux.Handle("DELETE /notes/{id}", mustHaveMiddlewares(http.HandlerFunc(server.DeleteNoteHandler)))

	log.Println("Запуск сервера 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Ошибка:", err)
	}
}
