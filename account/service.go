package account

import (
	"errors"
	"fmt"
	"math/rand"
	. "notes-app/entity"
	"notes-app/helper"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(input RegisterAccountInput) (Account, error)
	Login(input LoginAccountInput) (Account, error)
	GetAccountById(id int) (Account, error)
	VerifyEmail(token string) (bool, error)
	ResetPassword(input ResetPasswordInput) (string, string, error)
	VerifyTokenReset(token string) (bool, error)
	CreateAccount(input CreateAccountInput) (Account, error)
	//UpdateAccount(uri GetAccountUriInput, input UpdateAccountInput) (Account, error)
	DeleteAccount(input GetAccountUriInput) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) Register(input RegisterAccountInput) (Account, error) {
	account := Account{
		Email:       input.Email,
		FullName:    input.FullName,
		Role:        User,
		VerifyToken: helper.Encrypt(input.Email, os.Getenv("SALT_TOKEN_VERIFY")),
		ResetToken:  helper.Encrypt(input.Email, os.Getenv("SALT_TOKEN_VERIFY")),
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return account, err
	}

	account.Password = string(passwordHash)

	checkAccount, err := s.repository.FindByEmail(account.Email)

	if err != nil {
		return checkAccount, err
	}

	if checkAccount.ID != 0 {
		return checkAccount, errors.New("Email is Exist at the Database!")
	}

	newAccount, err := s.repository.Save(account)

	if err != nil {
		return newAccount, err
	}

	return newAccount, nil

}

func (s *service) Login(input LoginAccountInput) (Account, error) {
	account, err := s.repository.FindByEmail(input.Email)

	if err != nil {
		return account, err
	}

	if account.ID == 0 {
		return account, errors.New("no user found on that email")
	}

	if !account.IsVerified {
		return account, errors.New("Please kindly check your email to verify your account !")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password))

	if err != nil {
		return account, err
	}

	return account, nil

}

func (s *service) GetAccountById(id int) (Account, error) {
	account, err := s.repository.FindById(id)
	if err != nil {
		return account, err
	}
	return account, nil
}

func (s *service) VerifyEmail(token string) (bool, error) {

	account, err := s.repository.FindByTokenVerify(token)

	if err != nil {
		return false, err
	}

	if account.ID == 0 {
		return false, errors.New("no user found with that token")
	}

	account.IsVerified = true

	_, err = s.repository.Update(account)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *service) ResetPassword(input ResetPasswordInput) (string, string, error) {

	account, err := s.repository.FindByEmail(input.Email)

	if err != nil {
		return "", "", err
	}

	if account.ID == 0 {
		return "", "", errors.New("no user found with that email")
	}

	passwd := fmt.Sprintf("%v", rand.Int())

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)

	if err != nil {
		return "", "", err
	}

	account.Password = string(passwordHash)
	account.IsVerified = false

	updatedAccount, err := s.repository.Update(account)

	if err != nil {
		return "", "", err
	}

	return passwd, updatedAccount.ResetToken, nil

}

func (s *service) VerifyTokenReset(token string) (bool, error) {

	account, err := s.repository.FindByTokenReset(token)

	if err != nil {
		return false, err
	}

	if account.ID == 0 {
		return false, errors.New("no user found with that token")
	}

	account.IsVerified = true

	_, err = s.repository.Update(account)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *service) CreateAccount(input CreateAccountInput) (Account, error) {
	account := Account{
		Email:       input.Email,
		FullName:    input.FullName,
		Role:        UserRole(input.Role),
		VerifyToken: helper.Encrypt(input.Email, os.Getenv("SALT_TOKEN_VERIFY")),
		ResetToken:  helper.Encrypt(input.Email, os.Getenv("SALT_TOKEN_VERIFY")),
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return account, err
	}

	resp, err := helper.UploadImage(input.Image)

	if err != nil {
		return account, err
	}

	account.ImageUrl = resp.SecureURL

	account.Password = string(passwordHash)

	newAccount, err := s.repository.Save(account)

	if err != nil {
		return newAccount, err
	}

	return newAccount, nil

}

/*func (s *service) UpdateAccount(uri GetAccountUriInput, input UpdateAccountInput) (Account, error) {
	account, err := s.repository.FindById(uri.ID)

	if err != nil {
		return account, err
	}

}*/

func (s *service) DeleteAccount(uri GetAccountUriInput) (bool, error) {
	account, err := s.repository.FindById(uri.ID)

	if err != nil {
		return false, err
	}

	isDelete, err := s.repository.Delete(account)

	if err != nil {
		return isDelete, err
	}

	return isDelete, nil

}
