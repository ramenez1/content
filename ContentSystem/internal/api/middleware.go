package api

import (
	tools1 "ContentSystem/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

const SessionKey = "session_id"

type SessionAuth struct {
	rdb *redis.Client
}

func NewSessionAuth() *SessionAuth {
	s := &SessionAuth{}
	connRdb(s)
	return s
}

func (s *SessionAuth) Auth(ctx *gin.Context) {
	sessionID := ctx.GetHeader(SessionKey)
	//TODO : imp auth
	if sessionID == "" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, "session is empty")
	}
	//鉴权
	authKey := tools1.GetAuthKey(sessionID)
	loginTime, err := s.rdb.Get(ctx, authKey).Result()
	if err != nil && err != redis.Nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, "session auth error")
	}
	if loginTime == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "session auth fail")
	}

	fmt.Println("session id ", sessionID)
	ctx.Next()

}

func connRdb(s *SessionAuth) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("session_client error: ", err)
		panic(err)

	}
	s.rdb = rdb
}
