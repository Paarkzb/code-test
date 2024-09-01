package model

import "errors"

type Note struct {
	Id       int    `json:"-"`
	RfUserId int    `json:"rf_user_id" binding:"required"`
	Body     string `json:"body"`
}

type NoteResponse struct {
	Id   int          `json:"id"`
	User UserResponse `json:"user"`
	Body string       `json:"body"`
}

type UpdateNoteInput struct {
	Body *string `json:"body"`
}

func (i UpdateNoteInput) Validate() error {
	if i.Body == nil {
		return errors.New("поле body пустое")
	}

	return nil
}

type SpellerResponse struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}
