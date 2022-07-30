package itmoTrainerApi

import (
	"encoding/json"
	"net/http"
)

type problem struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Theme    string `json:"theme" db:"theme"`
	Rating   string `json:"rating" db:"rating"`
	Answers  string `json:"answers" db:"answers"`
	Text     string `json:"text" db:"text"`
	IsPublic string `json:"is_public" db:"isPublic"`
}

func problemExists(problemId string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM problems WHERE `id` = ?", problemId)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func (p *problem) load(problemId string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Get(p, "SELECT * FROM problems WHERE `id` = ?", problemId)
	if err != nil {
		return err
	}
	return nil
}

func CreateProblem(newProblem *problem) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	query := "INSERT INTO `problems`(`name`, `theme`, `rating`, `answers`, `text`, `isPublic`) VALUES (?,?,?,?,?,?)"
	_, err = db.Exec(query, newProblem.Name, newProblem.Theme, newProblem.Rating, newProblem.Answers, newProblem.Text, newProblem.IsPublic)
	if err != nil {
		return internalError(err)
	}
	return APIGatewayResponse{StatusCode: http.StatusOK}
}

func GetProblem(problemId string) APIGatewayResponse {
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	var p problem
	if err := p.load(problemId); err != nil {
		return internalError(err)
	}
	if res, err := json.Marshal(p); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}

func EditProblem(problemId string, newProblem *problem) APIGatewayResponse {
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	query := "UPDATE `problems` SET `name`=?,`theme`=?,`rating`=?,`answers`=?,`text`=?, `isPublic`=? WHERE id=?"
	_, err = db.Exec(query, newProblem.Name, newProblem.Theme, newProblem.Rating, newProblem.Answers, newProblem.Text, newProblem.IsPublic, newProblem.Id)
	if err != nil {
		return internalError(err)
	}
	return APIGatewayResponse{StatusCode: http.StatusOK}
}

func DeleteProblem(problemId string) APIGatewayResponse {
	if exists, err := problemExists(problemId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	query := "DELETE FROM `problems` WHERE `id` = ?"
	_, err = db.Exec(query, problemId)
	if err != nil {
		return internalError(err)
	}
	return APIGatewayResponse{StatusCode: http.StatusOK}
}
