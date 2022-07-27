package itmoTrainerApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type message struct {
	Id        string `json:"id" db:"id"`
	FromId    string `json:"from_id" db:"fromId"`
	ToId      string `json:"to_id" db:"toId"`
	ContestId string `json:"contest_id" db:"contestId"`
	Text      string `json:"text" db:"text"`
	Time      string `json:"time" db:"time"`
	Checked   bool   `json:"checked" db:"checked"`
}

func (m *message) load(id string) error {
	db, err := getConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Get(m, "SELECT * FROM contestMessages WHERE `id` = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func messageExists(id string) (bool, error) {
	db, err := getConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()
	cnt := 0
	err = db.Get(&cnt, "SELECT COUNT(`id`) FROM contestMessages WHERE `id` = ?", id)
	if err != nil {
		return false, err
	}
	return cnt == 1, nil
}

func GetMessagesFromContest(apiKey, contestId, userId string, adminMode bool) APIGatewayResponse {
	// checking empty parameters
	needCheckParams := map[string]string{
		"apiKey":    apiKey,
		"userId":    userId,
		"contestId": contestId,
	}
	for k, v := range needCheckParams {
		if len(v) == 0 {
			return APIGatewayResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Parameter " + k + " is required",
			}
		}
	}

	// checking is user allow to view messages
	v, err := userExists(userId)
	if err != nil {
		return internalError(err)
	}
	if !v {
		return APIGatewayResponse{StatusCode: http.StatusBadRequest, Body: fmt.Sprintf("User %s is not exists", userId)}
	}
	var u user
	err = u.load(userId)
	if err != nil {
		return internalError(err)
	}
	if (adminMode && apiKey != getAPIKEY()) || (!adminMode && u.Password != apiKey) {
		return accessDenied()
	}
	if v, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !v {
		return notFoundCustomText(fmt.Sprintf("Contest %s is not exists", contestId))
	}

	// connecting to database
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()

	// getting messages from database
	var messages []message
	users := make(map[string]user)
	if adminMode {
		err = db.Select(&messages, "SELECT * FROM contestMessages WHERE `contestId`=?", contestId)
	} else {
		err = db.Select(&messages, "SELECT * FROM contestMessages WHERE `contestId`=? AND (`toId`=? OR `fromId`=?)", contestId, userId, userId)
	}
	if err != nil {
		return internalError(err)
	}
	messagesResponse := make(map[string][]message)
	for _, i := range messages {
		for _, userId := range []string{i.ToId, i.FromId} {
			if _, ok := users[userId]; !ok {
				var u user
				err = u.load(userId)
				if err != nil {
					users[userId] = user{}
				} else {
					users[userId] = u
				}
			}
		}
		key := minString(i.ToId, i.FromId) + "_" + maxString(i.ToId, i.FromId)
		messagesResponse[key] = append(messagesResponse[key], i)
	}
	if len(messages) == 0 {
		return notFound()
	}
	body, err := json.Marshal(map[string]interface{}{
		"response": messagesResponse,
		"users":    users,
	})
	if err != nil {
		return internalError(err)
	}
	return APIGatewayResponse{StatusCode: http.StatusOK, Body: string(body)}
}

func CountUnreadMessages(apiKey, contestId, userId string, adminMode bool) APIGatewayResponse {
	// checking empty parameters
	needCheckParams := map[string]string{
		"apiKey":    apiKey,
		"contestId": contestId,
		"userId":    userId,
	}
	for k, v := range needCheckParams {
		if len(v) == 0 {
			return APIGatewayResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Parameter " + k + " is required",
			}
		}
	}

	// checking is userId exists
	if v, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !v && !adminMode {
		return notFoundCustomText(fmt.Sprintf("User %s not found", userId))
	}
	var u user
	err := u.load(userId)
	if err != nil {
		return internalError(err)
	}

	// checking is contestId exists
	if v, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !v {
		return notFoundCustomText(fmt.Sprintf("Contest %s not found", contestId))
	}

	// checking access to use method
	if (adminMode && apiKey != getAPIKEY()) || (!adminMode && u.Password != apiKey) {
		return accessDenied()
	}

	// connecting to database
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()

	// getting count of unread messages for each mode
	var res string
	if adminMode {
		if err := db.Get(&res, "SELECT COUNT(`id`) FROM `contestMessages` WHERE (`checked` = 'false' AND `contestId`=?)", contestId); err != nil {
			return internalError(err)
		}
	} else {
		if err := db.Get(&res, "SELECT COUNT(`id`) FROM `contestMessages` WHERE (`checked` = 'false' AND `toId`=? AND `contestId`=?)", userId); err != nil {
			return internalError(err)
		}
	}
	return APIGatewayResponse{StatusCode: http.StatusOK, Body: "{\"count\":\"" + res + "\"}"}
}

func sendMessage(apiKey, contestId, fromId, toId, text string, adminMode bool) APIGatewayResponse {
	// checking empty parameters
	needCheckParams := map[string]string{
		"apiKey":    apiKey,
		"fromId":    fromId,
		"toId":      toId,
		"contestId": contestId,
		"text":      text,
	}
	for k, v := range needCheckParams {
		if len(v) == 0 {
			return APIGatewayResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Parameter " + k + " is required",
			}
		}
	}
	// creating message
	msg := message{
		FromId:    fromId,
		ToId:      toId,
		ContestId: contestId,
		Text:      text,
		Time:      time.Now().Format("2006-01-02 15:04:05"),
		Checked:   false,
	}

	// checking is user allowed to use method
	var u user
	// checking is users exist
	for _, userId := range []string{toId, fromId} {
		exists, err := userExists(userId)
		if err != nil {
			return internalError(err)
		}
		if !exists {
			return notFoundCustomText(fmt.Sprintf("User %s is not exists", userId))
		}
	}
	err := u.load(fromId)
	if err != nil {
		return internalError(err)
	}
	if (adminMode && apiKey != getAPIKEY()) || (!adminMode && u.Password != apiKey) {
		return accessDenied()
	}

	// checking is contestId exists
	if v, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !v {
		return notFoundCustomText(fmt.Sprintf("Contest %s not found", contestId))
	}

	// connecting to database
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO `contestMessages`(`fromId`, `toId`, `contestId`, `text`, `time`, `checked`) VALUES (?, ?, ?, ?, ?, ?)",
		msg.FromId, msg.ToId, msg.ContestId, msg.Text, msg.Time, msg.Checked)
	if err != nil {
		return internalError(err)
	}
	return ok()
}

func CheckMessage(apiKey, userId, messageId string, adminMode bool) APIGatewayResponse {
	// checking empty parameters
	needCheckParams := map[string]string{
		"apiKey":    apiKey,
		"userId":    userId,
		"messageId": messageId,
	}
	for k, v := range needCheckParams {
		if len(v) == 0 {
			return APIGatewayResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Parameter " + k + " is required",
			}
		}
	}
	if v, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !v {
		return notFoundCustomText(fmt.Sprintf("User %s is not exists", userId))
	}
	var u user
	if err := u.load(userId); err != nil {
		return internalError(err)
	}
	if (adminMode && apiKey != getAPIKEY()) || (!adminMode && u.Password != apiKey) {
		return accessDenied()
	}

	if exists, err := messageExists(messageId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText(fmt.Sprintf("Message %s is not exists", messageId))
	}
	var m message
	if err := m.load(messageId); err != nil {
		return internalError(err)
	}
	if m.ToId != userId && !adminMode {
		return accessDenied()
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	_, err = db.Exec("UPDATE `contestMessages` SET `checked`='1' WHERE `id`=?", messageId)
	if err != nil {
		return internalError(err)
	}
	return ok()
}

func CheckAllChat(apiKey, contestId, userId string, adminMode bool) APIGatewayResponse {
	// checking empty parameters
	needCheckParams := map[string]string{
		"apiKey":    apiKey,
		"contestId": contestId,
		"userId":    userId,
	}
	for k, v := range needCheckParams {
		if len(v) == 0 {
			return APIGatewayResponse{
				StatusCode: http.StatusBadRequest,
				Body:       "Parameter " + k + " is required",
			}
		}
	}

	// checking is user allowed to edit message
	if v, err := userExists(userId); err != nil {
		return internalError(err)
	} else if !v {
		return notFoundCustomText(fmt.Sprintf("User %s is not exists", userId))
	}
	var u user
	if err := u.load(userId); err != nil {
		return internalError(err)
	}
	if (adminMode && apiKey != getAPIKEY()) || (!adminMode && u.Password != apiKey) {
		return accessDenied()
	}

	if exists, err := contestExists(contestId); err != nil {
		return internalError(err)
	} else if !exists {
		return notFoundCustomText(fmt.Sprintf("Contest %s is not exists", contestId))
	}
	db, err := getConnection()
	if err != nil {
		return internalError(err)
	}
	_, err = db.Exec("UPDATE `contestMessages` SET `checked`='1' WHERE `contestId`=? AND `toId`=?", contestId, userId)
	if err != nil {
		return internalError(err)
	}
	return ok()
}
