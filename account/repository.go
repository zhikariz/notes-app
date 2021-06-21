package account

import (
	. "notes-app/entity"

	"gorm.io/gorm"
)

type Repository interface {
	Save(account Account) (Account, error)
	FindByEmail(email string) (Account, error)
	FindByTokenVerify(token string) (Account, error)
	FindByTokenReset(token string) (Account, error)
	FindById(id int) (Account, error)
	Update(account Account) (Account, error)
	Delete(account Account) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(account Account) (Account, error) {
	err := r.db.Create(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repository) FindByEmail(email string) (Account, error) {
	var account Account
	err := r.db.Where("email = ?", email).Find(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repository) FindById(id int) (Account, error) {
	var account Account
	err := r.db.Where("id = ?", id).Find(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repository) Update(account Account) (Account, error) {
	err := r.db.Save(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repository) FindByTokenVerify(token string) (Account, error) {
	var account Account
	err := r.db.Where("verify_token = ?", token).Find(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repository) FindByTokenReset(token string) (Account, error) {
	var account Account
	err := r.db.Where("reset_token = ?", token).Find(&account).Error

	if err != nil {
		return account, err
	}

	return account, nil
}

func (r *repository) Delete(account Account) (bool, error) {
	err := r.db.Delete(&account).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
