package note

type CreateNotesInput struct {
	Title     string `json:"title" binding:"required"`
	Body      string `json:"body" binding:"required"`
	Secret    string `json:"secret,omitempty"`
	Type      string `json:"type" binding:"required"`
	AccountID int
}

type UpdateNotesInput struct {
	Title     string `json:"title" binding:"required"`
	Body      string `json:"body" binding:"required"`
	Secret    string `json:"secret,omitempty"`
	Type      string `json:"type" binding:"required"`
	AccountID int
}

type DeleteNotesInput struct {
	Secret string `json:"secret"`
}

type GetNotesUriInput struct {
	ID        int `uri:"id" binding:"required"`
	AccountID uint
	Secret    string
}

type GetNotesBodyInput struct {
	Secret string `json:"secret"`
}
