package itmo_trainer_api

import "testing"

func TestGetRecentAttempt(t *testing.T) {
	println(GetAttemptsListForUser("1", "-10").Body)
}
