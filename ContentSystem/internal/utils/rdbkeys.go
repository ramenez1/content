package utils

import "fmt"

func GetAuthKey(sessionID string) string {
	authKey := fmt.Sprintf("session_auth:%s", sessionID)
	return authKey
}

func GetSessionKey(userID string) string {
	sessionKey := fmt.Sprintf("session_id:%s", userID)
	return sessionKey
}
