package itmoTrainerApi

import "net/http"

type attempt struct {
	UserId     string `json:"user_id" db:"userId"`
	ContestId  string `json:"contest_id" db:"contestId"`
	Answers    string `json:"answers" db:"answers"`
	Result     string `json:"result" db:"result"`
	IsFinished bool   `json:"is_finished" db:"isFinished"`
	Id         string `json:"id" db:"id"`
}

func attemptExists(attemptId string) bool {
	return true
}

func (a *attempt) load(attemptId string) error {
	return nil
}

func GetAttempt(attemptId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func EditAttempt(attemptId string, newAttempt *attempt) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func SaveAttempt(attemptId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func CheckAttempt(attemptId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}
