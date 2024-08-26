package controllers

import (
	"database/sql"
	"myapp/internal/notes"
	"myapp/pkg/io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *HandlerLayer) CreateNote(w http.ResponseWriter, r *http.Request) {
	var text notes.NoteText
	err := io.ReadJSON(r, &text)
	if err != nil {
		io.SendError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	token := r.Context().Value("token").(string)
	note, err := h.LogicLayer.CreateNote(text.Text, token)
	if err != nil {
		io.SendError(w, "Failed to create note", http.StatusInternalServerError)
		h.Logger.Errorf("Error in create note: %v", err)
		return
	}
	h.Logger.Infof("Create Note: %s", note.ID)
	io.WriteJSON(w, http.StatusCreated, note)
}

func (h *HandlerLayer) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		io.SendError(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	var text notes.NoteText
	err = io.ReadJSON(r, &text)
	if err != nil {
		io.SendError(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	token := r.Context().Value("token").(string)
	note, err := h.LogicLayer.UpdateNote(id, text.Text, token)
	if err != nil {
		switch err.Error() {
		case "no permission to update note":
			io.SendError(w, "No permission to update note", http.StatusBadRequest)
			return
		case sql.ErrNoRows.Error():
			io.SendError(w, "Invalid note ID", http.StatusBadRequest)
			return
		}
		io.SendError(w, "Failed to update note", http.StatusInternalServerError)
		h.Logger.Errorf("Error in update note: %v", err)
		return
	}
	h.Logger.Infof("Update Note: %s", note.ID)
	io.WriteJSON(w, http.StatusOK, note)
}

func (h *HandlerLayer) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		io.SendError(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	token := r.Context().Value("token").(string)
	err = h.LogicLayer.DeleteNote(id, token)
	if err != nil {
		switch err.Error() {
		case "no permission to delete note":
			io.SendError(w, "No permission to delete note", http.StatusBadRequest)
			return
		case sql.ErrNoRows.Error():
			io.SendError(w, "Invalid note ID", http.StatusBadRequest)
			return
		}
		io.SendError(w, "Failed to delete note", http.StatusInternalServerError)
		h.Logger.Errorf("Error in delete note: %v", err)
		return
	}
	h.Logger.Infof("Delete Note: %s", id)
	io.SendError(w, "Deleted sucessfully", http.StatusOK)
}

func (h *HandlerLayer) GetNotes(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(string)
	notes, username, err := h.LogicLayer.GetNotes(token)
	if err != nil {
		io.SendError(w, "Failed to get notes", http.StatusInternalServerError)
		h.Logger.Errorf("Error in get all notes: %v", err)
		return
	}
	h.Logger.Infof("Get All Notes: %s", username)
	io.WriteJSON(w, http.StatusOK, notes)
}
