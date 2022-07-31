package itmo_trainer_api

import (
	"encoding/json"
	"errors"
	"net/http"
)

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

func contestExists(contestId string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM contests WHERE `id` = ?", contestId)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func (c *contest) load(contestId string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Get(c, "SELECT * FROM contests WHERE `id` = ?", contestId)
	if err != nil {
		return err
	}
	return nil
}

func (c *contest) getParticipants() (map[string]string, error) {
	res := make(map[string]string)
	err := json.Unmarshal([]byte(c.Participants), &res)
	return res, err
}

func GetContest(contestId string) APIGatewayResponse {
	if exists, err := contestExists(contestId); err != nil {
		return internalError(errors.New("Line 48: " + err.Error()))
	} else if !exists {
		return notFound()
	}
	var c contest
	if err := c.load(contestId); err != nil {
		return internalError(errors.New("Line 54: " + err.Error()))
	}
	if res, err := json.MarshalIndent(c, "", "    "); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}

}

func GetContestList(showHidden bool) APIGatewayResponse {
	var contestList []contest
	var query string
	if showHidden {
		query = "SELECT * FROM `contests` WHERE 1"
	} else {
		query = "SELECT * FROM `contests` WHERE `isPublic` = true"
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	rows, err := db.Queryx(query)
	for rows.Next() {
		var c contest
		err = rows.StructScan(&c)
		contestList = append(contestList, c)
	}
	if res, err := json.MarshalIndent(contestList, "", "    "); err != nil {
		return internalError(err)
	} else {
		return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(res)}
	}
}

func EditContest(contestId string, newContest *contest) APIGatewayResponse {
	if exists, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	query := "UPDATE `contests` SET `problemset`=?,`start`=?,`finish`=?,`disableTheory`=?,`isPublic`=?, `participants`=?, `name`=? WHERE id=?"
	_, err = db.Exec(query, newContest.Problemset, newContest.Start, newContest.Finish, newContest.DisableTheory,
		newContest.IsPublic, newContest.Participants, newContest.Name, contestId)
	if err != nil {
		return internalError(err)
	}
	return ok()
}

func CreateContest(newContest *contest) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	var problemset []string
	err = json.Unmarshal([]byte(newContest.Problemset), &problemset)
	if err != nil {
		return internalError(err)
	}
	for _, val := range problemset {
		if exists, err := problemExists(val); err != nil {
			return internalError(err)
		} else if !exists {
			return notFoundCustomText("Problem with ID " + val + " not found")
		}
	}
	query := "INSERT INTO `contests`(`problemset`, `start`, `finish`, `disableTheory`, `isPublic`, `participants`, `name`) VALUES (?,?,?,?,?,?,?)"
	_, err = db.Exec(query, newContest.Problemset, newContest.Start, newContest.Finish, newContest.DisableTheory,
		newContest.IsPublic, newContest.Participants, newContest.Name)
	if err != nil {
		return internalError(err)
	}
	return ok()
}

func DeleteContest(contestId string) APIGatewayResponse {
	if exists, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFound()
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	query := "DELETE FROM `contests` WHERE `id` = ?"
	_, err = db.Exec(query, contestId)
	if err != nil {
		return internalError(err)
	}
	return ok()
}

func ChangeParticipantState(contestId, userId, state string) APIGatewayResponse {
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()
	var c contest
	if exists, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText("Contest with ID " + contestId + " not found")
	}
	if err := c.load(contestId); err != nil {
		return internalError(err)
	}
	participants, err := c.getParticipants()
	if err != nil {
		return internalError(err)
	}
	if _, ok := participants[userId]; !ok {
		return notFoundCustomText("User with ID " + userId + " not found")
	} else {
		participants[userId] = state
	}
	if res, err := json.Marshal(participants); err != nil {
		return internalError(err)
	} else {
		query := "UPDATE `contests` SET `participants`=? WHERE `id`=?"
		_, err := db.Exec(query, res, contestId)
		if err != nil {
			return internalError(err)
		}
		return ok()
	}
}
