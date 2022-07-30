package itmoTrainerApi

import "net/http"

type note struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Text     string `json:"text" db:"text"`
	Category string `json:"category" db:"category"`
	Hide     bool   `json:"hide" db:"hide"`
	Date     string `json:"date" db:"date"`
}

func noteExists(noteId string) bool {
	return true
}

func (n *note) load(noteID string) error {
	return nil
}

func CanShowNotes() APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func GetNote(noteId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func GetNoteList(noteId string, showHidden bool) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func EditNote(noteId string, newNote *note) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func CreateNote(noteId string, newNote *note) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func DeleteNode(noteId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}
