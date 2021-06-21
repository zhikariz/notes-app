package note

import . "notes-app/entity"

type NoteFormatter struct {
	ID     uint     `json:"id"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Secret string   `json:"secret"`
	Type   NoteType `json:"type" `
}

func FormatNote(note Note) (noteFormatter NoteFormatter) {
	noteFormatter = NoteFormatter{
		ID:     note.ID,
		Title:  note.Title,
		Body:   note.Body,
		Secret: note.Secret,
		Type:   note.Type,
	}
	return
}

func FormatNotes(notes []Note) (notesFormatter []NoteFormatter) {
	if len(notes) == 0 {
		return []NoteFormatter{}
	}

	for _, note := range notes {
		formatter := FormatNote(note)
		notesFormatter = append(notesFormatter, formatter)
	}
	return
}
