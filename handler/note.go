package handler

import (
	"notes-app/auth"
	. "notes-app/entity"
	"notes-app/helper"
	"notes-app/note"

	"github.com/gin-gonic/gin"
)

type noteHandler struct {
	authService auth.Service
	noteService note.Service
}

func NewNoteHandler(authService auth.Service, noteService note.Service) *noteHandler {
	return &noteHandler{authService, noteService}
}

func (h *noteHandler) GetAllNote(c *gin.Context) {
	notes, err := h.noteService.GetNotes()

	if err != nil {
		helper.ErrorHandling(c, err, "Error Displaying All Note")
		return
	}
	formatter := note.FormatNotes(notes)
	helper.SuccessHandling(c, formatter, "Successfully Displaying All Note")

}

func (h *noteHandler) GetNoteByAccountId(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(Account)

	noteAccount, err := h.noteService.GetNoteByAccountId(int(currentUser.ID))

	if err != nil {
		helper.ErrorHandling(c, err, "Error Displaying Note")
		return
	}

	formatter := note.FormatNotes(noteAccount)
	helper.SuccessHandling(c, formatter, "Successfully Displaying Note Account")
}

func (h *noteHandler) GetNoteById(c *gin.Context) {
	var uri note.GetNotesUriInput
	err := c.ShouldBindUri(&uri)

	if err != nil {
		helper.ErrorValidation(c, err, "Error Displaying Note")
		return
	}

	secret := c.Query("secret")
	uri.Secret = secret

	existNote, err := h.noteService.GetNote(uri)

	if err != nil {
		helper.ErrorHandling(c, err, "Error Displaying Note")
		return
	}

	formatter := note.FormatNote(existNote)
	helper.SuccessHandling(c, formatter, "Successfully Displaying Note")
}

func (h *noteHandler) CreateNote(c *gin.Context) {
	var input note.CreateNotesInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		helper.ErrorValidation(c, err, "Error Create Note")
		return
	}

	currentUser := c.MustGet("currentUser").(Account)

	input.AccountID = int(currentUser.ID)

	newNote, err := h.noteService.CreateNote(input)

	if err != nil {
		helper.ErrorHandling(c, err, "Error Create Note")
		return
	}

	formatter := note.FormatNote(newNote)
	helper.SuccessHandling(c, formatter, "Successfuly Create Note")
}

func (h *noteHandler) UpdateNote(c *gin.Context) {
	var uri note.GetNotesUriInput
	var input note.UpdateNotesInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed Update Note")
		return
	}

	err = c.ShouldBindJSON(&input)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed Update Note")
		return
	}

	currentUser := c.MustGet("currentUser").(Account)

	secret := c.Query("secret")
	uri.Secret = secret

	input.AccountID = int(currentUser.ID)

	updatedNote, err := h.noteService.UpdateNote(uri, input)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed Update Note")
		return
	}

	formatter := note.FormatNote(updatedNote)

	helper.SuccessHandling(c, formatter, "Successfully Update Note")

}

func (h *noteHandler) DeleteNote(c *gin.Context) {
	var uri note.GetNotesUriInput

	err := c.ShouldBindUri(&uri)

	if err != nil {
		helper.ErrorValidation(c, err, "Failed to Delete Note")
		return
	}

	secret := c.Query("secret")
	uri.Secret = secret

	currentUser := c.MustGet("currentUser").(Account)
	uri.AccountID = currentUser.ID

	isDeleted, err := h.noteService.DeleteNote(uri)

	if err != nil {
		helper.ErrorHandling(c, err, "Failed to Delete Note")
		return
	}

	data := gin.H{"is_deleted": isDeleted}
	metaMessage := "Notes cannot be deleted !"

	if isDeleted {
		metaMessage = "Notes has been deleted !"
	}

	helper.SuccessHandling(c, data, metaMessage)

}
