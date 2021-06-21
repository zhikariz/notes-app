package note

import (
	. "notes-app/entity"

	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Note, error)
	FindById(id int) (Note, error)
	FindByAccountId(accountId int) ([]Note, error)
	Save(note Note) (Note, error)
	Update(note Note) (Note, error)
	Delete(note Note) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Note, error) {
	var notes []Note
	err := r.db.Find(&notes).Error

	if err != nil {
		return notes, err
	}
	return notes, nil
}

func (r *repository) FindById(id int) (Note, error) {
	var note Note
	err := r.db.Where("id = ?", id).Find(&note).Error

	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) FindByAccountId(accountId int) ([]Note, error) {
	var note []Note
	err := r.db.Where("account_id", accountId).Find(&note).Error

	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) Save(note Note) (Note, error) {
	err := r.db.Create(&note).Error

	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) Update(note Note) (Note, error) {
	err := r.db.Save(&note).Error

	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) Delete(note Note) (bool, error) {
	err := r.db.Delete(&note).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
