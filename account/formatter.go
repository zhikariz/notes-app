package account

import . "notes-app/entity"

type AccountFormatter struct {
	ID         uint     `json:"id"`
	Email      string   `json:"email"`
	FullName   string   `json:"full_name"`
	Role       UserRole `json:"role"`
	ImageUrl   string   `json:"image_url"`
	Token      string   `json:"token"`
	IsVerified bool     `json:"is_verified"`
}

func FormatAccount(account Account, token string) (accountResponse AccountFormatter) {
	accountResponse = AccountFormatter{
		ID:         account.ID,
		Email:      account.Email,
		FullName:   account.FullName,
		Role:       account.Role,
		ImageUrl:   account.ImageUrl,
		Token:      token,
		IsVerified: account.IsVerified,
	}
	return
}
