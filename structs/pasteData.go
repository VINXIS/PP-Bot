package structs

import "strings"

// PasteData holds the paste data to send to paste.ee
type PasteData struct {
	sections []section
}

// Section is a section of the paste data
type section struct {
	name     string
	contents string
}

// PasteResult holds the paste data sent from paste.ee
type PasteResult struct {
	ID      string `json:"id"`
	Link    string `json:"link"`
	Success bool   `json:"success"`
}

// NewPasteData creates a new paste data to send to paste.ee
func NewPasteData(username, content string) PasteData {
	return PasteData{
		sections: []section{
			section{
				name:     "PP Changes for " + username,
				contents: content,
			},
		},
	}
}

// Marshal returns the data into []byte
func (d *PasteData) Marshal() []byte {
	text := `{"sections":[{"name":"` + d.sections[0].name + `","contents":"` + strings.Replace(strings.Replace(strings.Replace(d.sections[0].contents, "\n", "\\n", -1), "\r", "\\r", -1), "\"", "\\\"", -1) + `"}]}`
	return []byte(text)
}
