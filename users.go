package itmoTrainerApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type user struct {
	Email        string `json:"email" db:"email"`
	Password     string `json:"password" db:"password"`
	Name         string `json:"name" db:"name"`
	Surname      string `json:"surname" db:"surname"`
	Patronymic   string `json:"patronymic" db:"patronymic"`
	Login        string `json:"login" db:"login"`
	Birthday     string `json:"birthday" db:"birthday"`
	IsActive     bool   `json:"is_active" db:"isActive"`
	Id           string `json:"id" db:"id"`
	Rating       int    `json:"rating" db:"rating"`
	IsAdmin      bool   `json:"is_admin" db:"isAdmin"`
	Course       int    `json:"course" db:"course"`
	RegisterHash string `json:"register_hash" db:"registerHash"`
	TimeCreated  string `json:"time_created" db:"timeCreated"`
	Avatar       string `json:"avatar" db:"avatar"`
}

func (u *user) getFullName() string {
	return strings.Trim(u.Surname+" "+u.Name+" "+u.Patronymic, " ")
}

func userExists(id string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM users WHERE `id` = ?", id)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func (u *user) load(id string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Get(u, "SELECT * FROM users WHERE `id` = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(userId string) APIGatewayResponse {
	if len(userId) == 0 {
		return APIGatewayResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Parameter userId is empty",
		}
	}
	if v, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !v {
		return notFoundCustomText(fmt.Sprintf("User %s not found", userId))
	}
	var u user
	err := u.load(userId)
	if err != nil {
		return internalError(err)
	}
	u.Password = ""
	res, err := json.Marshal(u)
	if err != nil {
		return internalError(err)
	}
	return APIGatewayResponse{
		StatusCode: http.StatusOK,
		Body:       string(res),
	}
}

func UpdateUser(userId string, newUser *user) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func DeleteUser(userId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func CreateUser(userId string, newUser *user) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}

func ApproveUser(userId string) APIGatewayResponse {
	return APIGatewayResponse{StatusCode: http.StatusNotImplemented}
}
