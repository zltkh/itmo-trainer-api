package itmoTrainerApi

import "net/http"

type problem struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Theme    string `json:"theme" db:"theme"`
	Rating   string `json:"rating" db:"rating"`
	Answers  string `json:"answers" db:"answers"`
	Text     string `json:"text" db:"text"`
	IsPublic string `json:"is_public" db:"isPublic"`
}

func problemExists(problemId string) bool {
	return true
}

func (p *problem) load(problemId string) error {
	return nil
}

func GetProblem(problemId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func EditProblem(problemId string, newProblem *problem) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func DeleteProblem(problemId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}
