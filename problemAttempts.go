package itmo_trainer_api

import (
	"encoding/json"
	"net/http"
)

type problemAttempt struct {
	Id        string `json:"id" db:"id"`
	ProblemId string `json:"problem_id" db:"problemId"`
	Answer    string `json:"answer" db:"answer"`
	Verdict   bool   `json:"verdict" db:"verdict"`
	Time      string `json:"time" db:"time"`
	UserId    string `json:"user_id" db:"userId"`
	Source    string `json:"source" db:"source"`
}

func problemAttemptExists(problemAttemptId string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM problemsAttempts WHERE `id` = ?", problemAttemptId)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func (p *problemAttempt) load(problemAttemptId string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Get(p, "SELECT * FROM problemsAttempts WHERE `id` = ?", problemAttemptId)
	if err != nil {
		return err
	}
	return nil
}

func GetProblemAttempt(problemAttemptId string) APIGatewayResponse {
	if exists, err := problemExists(problemAttemptId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	var p problemAttempt
	if err := p.load(problemAttemptId); err != nil {
		return internalError(err)
	}
	if res, err := json.MarshalIndent(p, "", "    "); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}

func CreateProblemAttempt(source, problemId, userId, answer string) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	var p problem
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("Problem with ID " + problemId + " not found")
	}
	if exists, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("User with ID " + userId + " not found")
	}
	if err := p.load(problemId); err != nil {
		return internalError(err)
	}
	attempt, err := p.check(source, userId, answer)
	if err != nil {
		return internalError(err)
	}
	query := "INSERT INTO `problemsAttempts`(`problemId`, `answer`, `verdict`, `time`, `userId`, `source`) VALUES (?,?,?,?,?,?)"
	if _, err = db.Exec(query, attempt.ProblemId, attempt.Answer, attempt.Verdict, attempt.Time, attempt.UserId, attempt.Source); err != nil {
		return internalError(err)
	}
	return ok()
}

func GetAttemptsList(problemId string) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("Problem with ID " + problemId + " not found")
	}
	query := "SELECT * FROM `problemsAttempts` WHERE `problemId` = ?"
	var p []problemAttempt
	if err := db.Select(&p, query, problemId); err != nil {
		return internalError(err)
	}
	if res, err := json.MarshalIndent(p, "", "    "); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}

func GetAttemptsListForUser(problemId, userId string) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("Problem with ID " + problemId + " not found")
	}
	if exists, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("User with ID " + userId + " not found")
	}
	query := "SELECT * FROM `problemsAttempts` WHERE `problemId` = ? AND `userId` = ?"
	var p []problemAttempt
	if err := db.Select(&p, query, problemId, userId); err != nil {
		return internalError(err)
	}
	if res, err := json.MarshalIndent(p, "", "    "); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}

func GetRecentAttempt(problemId, userId string) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("Problem with ID " + problemId + " not found")
	}
	if exists, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("User with ID " + userId + " not found")
	}
	query := "SELECT * FROM `problemsAttempts` WHERE `problemId` = ? AND `userId` = ? ORDER BY `time` DESC LIMIT 1" // TODO change
	var p problemAttempt
	if err := db.Get(&p, query, problemId, userId); err != nil {
		return internalError(err)
	}
	if res, err := json.Marshal(p); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}
