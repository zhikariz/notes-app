package account

type RegisterAccountInput struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginAccountInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ResetPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyEmailInput struct {
	Token string `uri:"token" binding:"required"`
}

type GetAccountUriInput struct {
	ID   int `uri:"id" binding:"required"`
	Role string
}

type CreateAccountInput struct {
	Email    string `form:"email" binding:"required,email"`
	FullName string `form:"fullname" binding:"required"`
	Password string `form:"password" binding:"required"`
	Role     string `form:"role" binding:"required"`
	Image    string
}

type UpdateAccountInput struct {
	Email    string `form:"email" binding:"required,email"`
	FullName string `form:"fullname" binding:"required"`
	Password string `form:"password" binding:"required"`
	Role     string `form:"role" binding:"required"`
	Image    string
}
