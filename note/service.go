package note

import (
	"errors"
	. "notes-app/entity"
)

type Service interface {
	GetNotes() ([]Note, error)
	GetNote(uri GetNotesUriInput) (Note, error)
	GetNoteByAccountId(accountId int) ([]Note, error)
	CreateNote(input CreateNotesInput) (Note, error)
	UpdateNote(uri GetNotesUriInput, input UpdateNotesInput) (Note, error)
	DeleteNote(input GetNotesUriInput) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetNotes() ([]Note, error) {
	notes, err := s.repository.FindAll()

	if err != nil {
		return notes, err
	}

	return notes, nil
}

func (s *service) GetNote(uri GetNotesUriInput) (Note, error) {
	note, err := s.repository.FindById(uri.ID)

	if err != nil {
		return note, err
	}

	if note.ID == 0 {
		return note, errors.New("There is no such note with that id")
	}
	var isUsingSecret bool

	if note.Secret != "" {
		isUsingSecret = true
	}

	if isUsingSecret {
		if uri.Secret == "" {
			return note, errors.New("Need Secret for this note.")
		}

		if note.Secret != uri.Secret {
			return note, errors.New("Secret Notes is wrong.")
		}
	}

	return note, nil
}

func (s *service) GetNoteByAccountId(accountId int) ([]Note, error) {
	note, err := s.repository.FindByAccountId(accountId)

	if err != nil {
		return note, err
	}

	return note, nil
}

func (s *service) CreateNote(input CreateNotesInput) (Note, error) {
	note := Note{
		Title:     input.Title,
		Body:      input.Body,
		Secret:    input.Secret,
		Type:      NoteType(input.Type),
		AccountID: uint(input.AccountID),
	}

	note, err := s.repository.Save(note)

	if err != nil {
		return note, err
	}

	return note, nil
}

func (s *service) UpdateNote(uri GetNotesUriInput, input UpdateNotesInput) (Note, error) {
	note, err := s.repository.FindById(uri.ID)

	if err != nil {
		return note, err
	}

	if note.AccountID != uint(input.AccountID) {
		return note, errors.New("Not an owner of the notes")
	}

	var isUsingSecret bool

	if note.Secret != "" {
		isUsingSecret = true
	}

	if isUsingSecret {

		if uri.Secret == "" {
			return note, errors.New("Need Secret for this note.")
		}

		if note.Secret != uri.Secret {
			return note, errors.New("Secret Notes is wrong.")
		}
	}

	note.Title = input.Title
	note.Body = input.Body
	note.Secret = input.Secret
	note.Type = NoteType(input.Type)

	newUpdateNote, err := s.repository.Update(note)

	if err != nil {
		return newUpdateNote, err
	}

	return newUpdateNote, nil

}

func (s *service) DeleteNote(input GetNotesUriInput) (bool, error) {
	note, err := s.repository.FindById(input.ID)

	if err != nil {
		return false, err
	}

	if note.AccountID != input.AccountID {
		return false, errors.New("Not an owner of the notes")
	}

	var isUsingSecret bool

	if note.Secret != "" {
		isUsingSecret = true
	}

	if isUsingSecret {

		if input.Secret == "" {
			return false, errors.New("Need Secret for this note.")
		}

		if note.Secret != input.Secret {
			return false, errors.New("Secret Notes is wrong.")
		}
	}

	isDeleted, err := s.repository.Delete(note)

	if err != nil {
		return isDeleted, err
	}

	return isDeleted, nil
}
