package itmoTrainerApi

import "net/http"

type contest struct {
	Id            string `json:"id" db:"id"`
	Problemset    string `json:"problemset" db:"problemset"`
	Start         string `json:"start" db:"start"`
	Finish        string `json:"finish" db:"finish"`
	DisableTheory bool   `json:"disable_theory" db:"disableTheory"`
	IsPublic      bool   `json:"is_public" db:"isPublic"`
	Participants  string `json:"participants" db:"participants"`
	Name          string `json:"name" db:"name"`
}

func contestExists(id string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM contests WHERE `id` = ?", id)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func (c *contest) loadContest(contestId string) error {
	return nil
}

func GetContest(contestId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func GetContestList(apiKey string, showHidden bool) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func EditContest(contestId string, newContest *contest) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func CreateContest(contestId string, newContest *contest) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func DeleteContest(contestId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}
