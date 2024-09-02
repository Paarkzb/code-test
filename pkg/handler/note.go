package handler

import (
	"codetest/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func (h *Handler) getAllNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)

	if err != nil {
		return
	}

	notes, err := h.service.Note.GetAll(userId)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) saveNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		return
	}

	var input model.Note
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	noteId, err := h.service.Note.Create(userId, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": noteId})
}

func (h *Handler) getNoteById(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		return
	}

	paramId := chi.URLParam(r, "noteId")
	logrus.Println("paramId", paramId)
	id, err := strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "неправильный параметр id")
		return
	}

	note, err := h.service.Note.GetById(userId, id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		return
	}

	paramId := chi.URLParam(r, "noteId")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "неправильный параметр id")
		return
	}

	var input model.UpdateNoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.Note.Update(userId, id, input)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse{Status: "ok"})
}

func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(w, r)
	if err != nil {
		return
	}

	paramId := chi.URLParam(r, "noteId")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, "неправильный параметр id")
		return
	}

	err = h.service.Note.Delete(userId, id)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResponse{Status: "ok"})
}
