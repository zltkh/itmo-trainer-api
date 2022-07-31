package itmo_trainer_api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type contestAttempt struct {
	UserId     string `json:"user_id" db:"userId"`
	ContestId  string `json:"contest_id" db:"contestId"`
	Answers    string `json:"answers" db:"answers"`
	Result     string `json:"result" db:"result"`
	IsFinished bool   `json:"is_finished" db:"isFinished"`
	Id         string `json:"id" db:"id"`
}

func attemptExists(attemptId string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM contestAttempts WHERE `id` = ?", attemptId)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func (a *contestAttempt) load(attemptId string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Get(a, "SELECT * FROM contestAttempts WHERE `id` = ?", attemptId)
	if err != nil {
		return err
	}
	return nil
}

func (a *contestAttempt) getAnswers() (res map[string]string, err error) {
	db, err := getConnection()
	if err != nil {
		return
	}
	defer db.Close()
	err = json.Unmarshal([]byte(a.Answers), &res)
	return
}

func GetAttempt(attemptId string) APIGatewayResponse {
	if exists, err := attemptExists(attemptId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	var a contestAttempt
	if err := a.load(attemptId); err != nil {
		return internalError(err)
	}
	if res, err := json.Marshal(a); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}

func EditAttempt(attemptId string, newAttempt *contestAttempt) APIGatewayResponse {
	if exists, err := attemptExists(attemptId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	query := "UPDATE `contestAttempts` SET `userId`=?,`contestId`=?,`answers`=?,`result`=?,`isFinished`=? WHERE id=?"
	_, err = db.Exec(query, newAttempt.UserId, newAttempt.ContestId, newAttempt.Answers, newAttempt.Result, newAttempt.IsFinished)
	if err != nil {
		return internalError(err)
	}
	return ok()
}

func CheckAttempt(attemptId string, wg *sync.WaitGroup) APIGatewayResponse {
	defer wg.Done()
	type result struct {
		Verdict bool   `json:"verdict"`
		Rating  string `json:"rating"`
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	var a contestAttempt
	if exists, err := attemptExists(attemptId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	if err := a.load(attemptId); err != nil {
		return internalError(err)
	}
	var c contest
	if exists, err := contestExists(a.ContestId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	if err := c.load(a.ContestId); err != nil {
		return internalError(err)
	}
	answersUser, err := a.getAnswers()
	if err != nil {
		return internalError(err)
	}
	var problemset []string
	err = json.Unmarshal([]byte(c.Problemset), &problemset)
	if err != nil {
		return internalError(err)
	}
	query := "SELECT * FROM `problems` WHERE `id` = "
	for i, problemId := range problemset {
		query += "'" + problemId + "'"
		if i+1 != len(problemset) {
			query += " OR `id` = "
		}
	}
	var problems []problem
	if err = db.Select(&problems, query); err != nil {
		return internalError(err)
	}
	checked := make(map[string]result)
	for _, p := range problems {
		problemId := p.Id
		var answers []string
		err = json.Unmarshal([]byte(p.Answers), &answers)
		solved := false
		for _, ans := range answers {
			if ans == answersUser[problemId] {
				solved = true
				break
			}
		}
		checked[problemId] = result{
			Verdict: solved,
			Rating:  p.Rating,
		}
	}
	if res, err := json.MarshalIndent(checked, "", "    "); err != nil {
		return internalError(err)
	} else {
		query := "UPDATE `contestAttempts` SET `result`=?,`isFinished`=true WHERE `id`=?"
		_, err = db.Exec(query, string(res), attemptId)
		if err != nil {
			return internalError(err)
		}
		return ok()
	}
}

func CheckAllAttempts(contestId string) APIGatewayResponse {
	var attempts []string
	query := "SELECT `id` FROM `contestAttempts` WHERE `contestId`=?"
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	if err = db.Select(&attempts, query, contestId); err != nil {
		return internalError(err)
	}
	var wg sync.WaitGroup
	for _, attemptId := range attempts {
		wg.Add(1)
		go CheckAttempt(attemptId, &wg)
	}
	wg.Wait()
	return ok()
}
