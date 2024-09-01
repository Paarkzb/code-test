package service

import (
	"crypto/tls"
	"encoding/json"
	"medodstest/internal/model"
	"medodstest/pkg/repository"
	"net/http"
	"strings"
)

type NoteService struct {
	repo repository.Note
}

func NewNoteService(repo repository.Note) *NoteService {
	return &NoteService{
		repo: repo,
	}
}

func (s *NoteService) Create(userId int, note model.Note) (int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, "http://speller.yandex.net/services/spellservice.json/checkText", nil)
	if err != nil {
		return 0, err
	}

	q := req.URL.Query()
	q.Add("text", note.Body)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var sr []model.SpellerResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return 0, err
	}

	for _, v := range sr {
		note.Body = strings.Replace(note.Body, v.Word, v.S[0], 1)
	}

	return s.repo.Create(userId, note)
}

func (s *NoteService) GetAll(userId int) ([]model.NoteResponse, error) {
	return s.repo.GetAll(userId)
}

func (s *NoteService) GetById(userId, noteId int) (model.NoteResponse, error) {
	return s.repo.GetById(userId, noteId)
}

func (s *NoteService) Delete(userId, noteId int) error {
	return s.repo.Delete(userId, noteId)
}

func (s *NoteService) Update(userId, noteId int, input model.UpdateNoteInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, "http://speller.yandex.net/services/spellservice.json/checkText", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("text", *input.Body)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var sr []model.SpellerResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return err
	}

	for _, v := range sr {
		*input.Body = strings.Replace(*input.Body, v.Word, v.S[0], 1)
	}

	return s.repo.Update(userId, noteId, input)
}
