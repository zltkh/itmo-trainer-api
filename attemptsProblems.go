package itmoTrainerApi

import "net/http"

type problemAttempt struct {
	Id        string `json:"id" db:"id"`
	ProblemId string `json:"problem_id" db:"problemId"`
	Answer    string `json:"answer" db:"answer"`
	Verdict   bool   `json:"verdict" db:"verdict"`
	Time      string `json:"time" db:"time"`
	UserId    string `json:"user_id" db:"userId"`
	Source    string `json:"source" db:"source"`
}

func problemAttemptExists(problemAttemptId string) bool {
	return true
}

func (p *problemAttempt) load(problemAttemptId string) error {
	return nil
}

func GetProblemAttempt(problemAttemptId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func EditProblemAttempt(problemAttemptId string, newAttempt *attempt) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func CheckProblemAttempt(problemAttemptId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func GetAttemptsList(problemId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func GetAttemptsListForUser(problemId, userId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func GetRecentAttempt(problemId, userId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}
