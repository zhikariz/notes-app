package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"notes-app/account"
	"notes-app/auth"
	. "notes-app/entity"
	"notes-app/helper"
	"path"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

type accountHandler struct {
	authService    auth.Service
	accountService account.Service
}

func NewAccountHandler(authService auth.Service, accountService account.Service) *accountHandler {
	return &accountHandler{authService, accountService}
}

func (h *accountHandler) RegisterAccount(c *gin.Context) {
	var input account.RegisterAccountInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		helper.ErrorValidation(c, err, "Register Account Failed !")
		return
	}

	newUser, err := h.accountService.Register(input)

	if err != nil {
		helper.ErrorHandling(c, err, "Register Account Failed !")
		return
	}

	token, err := h.authService.GenerateToken(newUser)

	if err != nil {
		helper.ErrorHandling(c, err, "Register Account Failed !")
		return
	}

	tmpl, err := template.ParseFiles(path.Join("assets", "verify-email.html"))

	urls := url.URL{
		Host: location.Get(c).Host,
		Path: "api/v1/verify-account",
	}

	dataTemplate := map[string]interface{}{
		"email":     newUser.Email,
		"full_name": newUser.FullName,
		"url":       path.Join(urls.String(), newUser.VerifyToken),
		"scheme":    location.Get(c).Scheme,
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, dataTemplate); err != nil {
		helper.ErrorHandling(c, err, "Register Account Failed !")
		return
	}

	result := tpl.String()

	go helper.SendMail(newUser.Email, "Notes Account Successfully Registered", result)

	formatter := account.FormatAccount(newUser, token)
	helper.SuccessHandling(c, formatter, "Account successfully registered")
}

func (h *accountHandler) Login(c *gin.Context) {
	var input account.LoginAccountInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		helper.ErrorValidation(c, err, "Login Failed !")
		return
	}

	loggedInUser, err := h.accountService.Login(input)
	if err != nil {
		helper.ErrorHandling(c, err, "Login Failed !")
		return
	}

	token, err := h.authService.GenerateToken(loggedInUser)
	if err != nil {
		helper.ErrorHandling(c, err, "Login Failed !")
		return
	}

	formatter := account.FormatAccount(loggedInUser, token)
	helper.SuccessHandling(c, formatter, "Successfully Logged in !")
}

func (h *accountHandler) VerifyEmail(c *gin.Context) {
	var uri account.VerifyEmailInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed to Verify Token")
		return
	}

	isVerified, err := h.accountService.VerifyEmail(uri.Token)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Verify Token")
		return
	}

	data := gin.H{"is_verified": isVerified}
	metaMessage := "Account cannot be verified !"

	if isVerified {
		metaMessage = "Account has been verified !"
	}

	helper.SuccessHandling(c, data, metaMessage)
}

func (h *accountHandler) ResetPassword(c *gin.Context) {
	var input account.ResetPasswordInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed to Reset Password")
		return
	}

	newPassword, resetToken, err := h.accountService.ResetPassword(input)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Reset Password")
		return
	}

	tmpl, err := template.ParseFiles(path.Join("assets", "reset-password.html"))

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Reset Password")
		return
	}

	urls := url.URL{
		Host: location.Get(c).Host,
		Path: "api/v1/verify-password",
	}

	dataTemplate := map[string]interface{}{
		"email":       input.Email,
		"newPassword": newPassword,
		"url":         path.Join(urls.String(), resetToken),
		"scheme":      location.Get(c).Scheme,
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, dataTemplate); err != nil {
		helper.ErrorHandling(c, err, "Failed to Reset Password!")
		return
	}

	result := tpl.String()

	go helper.SendMail(input.Email, "Successfully Reset Password", result)

	data := gin.H{"is_reset": true, "message": "Kindly check your email, your password and verify it !"}
	metaMessage := "Account Password Has Been Reset!"

	helper.SuccessHandling(c, data, metaMessage)

}

func (h *accountHandler) VerifyReset(c *gin.Context) {
	var uri account.VerifyEmailInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed to Verify Token")
		return
	}

	isVerified, err := h.accountService.VerifyTokenReset(uri.Token)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Verify Token")
		return
	}

	data := gin.H{"is_verified": isVerified}
	metaMessage := "Account cannot be verified !"

	if isVerified {
		metaMessage = "Account has been verified !"
	}

	helper.SuccessHandling(c, data, metaMessage)
}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	var input account.CreateAccountInput

	err := c.ShouldBind(&input)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed to Create an Account")
		return
	}

	file, err := c.FormFile("avatar")

	if err != nil {
		helper.ErrorValidation(c, err, "Failed to Create an Account")
		return
	}

	currentUser := c.MustGet("currentUser").(Account)
	userId := currentUser.ID

	paths := fmt.Sprintf("images/%d-%s", userId, file.Filename)
	_ = c.SaveUploadedFile(file, paths)
	input.Image = paths

	newAccount, err := h.accountService.CreateAccount(input)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Create an Account")
		return
	}

	token, err := h.authService.GenerateToken(newAccount)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Create an Account!")
		return
	}

	tmpl, err := template.ParseFiles(path.Join("assets", "verify-email.html"))

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Create an Account")
		return
	}

	urls := url.URL{
		Host: location.Get(c).Host,
		Path: "api/v1/verify-account",
	}

	dataTemplate := map[string]interface{}{
		"email":     newAccount.Email,
		"full_name": newAccount.FullName,
		"url":       path.Join(urls.String(), newAccount.VerifyToken),
		"scheme":    location.Get(c).Scheme,
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, dataTemplate); err != nil {
		helper.ErrorHandling(c, err, "Failed to Create an Account !")
		return
	}

	result := tpl.String()

	go helper.SendMail(newAccount.Email, "Account Successfully Registered", result)

	formatter := account.FormatAccount(newAccount, token)
	helper.SuccessHandling(c, formatter, "Account successfully registered")

}
