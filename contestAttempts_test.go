package itmo_trainer_api

import "testing"

func TestAttempt_getAnswers(t *testing.T) {
	var a contestAttempt
	a.load("-1")
	res, err := a.getAnswers()
	if err != nil {
		println(err.Error())
	} else {
		for key, val := range res {
			println(key + ": " + val)
		}
	}
}

func TestCheckAttempt(t *testing.T) {
	res := CheckAllAttempts("1")
	println(res.Body)
}
