package service

import (
	"codetest/internal/model"
	"codetest/pkg/repository"
	"crypto/tls"
	"encoding/json"
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
	text, err := checkText(note.Body)
	if err != nil {
		return 0, err
	}
	note.Body = text

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

	text, err := checkText(*input.Body)
	if err != nil {
		return err
	}
	*input.Body = text

	return s.repo.Update(userId, noteId, input)
}

func checkText(text string) (string, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(http.MethodGet, "http://speller.yandex.net/services/spellservice.json/checkText", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var sr []model.SpellerResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", err
	}

	for _, v := range sr {
		text = strings.Replace(text, v.Word, v.S[0], 1)
	}

	return text, nil
}
