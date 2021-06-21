package entity

import "gorm.io/gorm"

type UserRole string

const (
	Admin = "Admin"
	User  = "User"
)

type Account struct {
	gorm.Model
	Email       string
	FullName    string
	Password    string
	ImageUrl    string
	IsVerified  bool
	VerifyToken string
	ResetToken  string
	Role        UserRole
}
