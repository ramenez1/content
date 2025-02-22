package service

import (
	"ContentSystem/internal/dao"
	"ContentSystem/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginReq struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
}

func (c *CmsApp) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accountDao := dao.NewAccountDao(c.db)
	account, err := accountDao.FirstByUserId(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "账号错误"})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "密码错误"})
		return
	}
	sessionID, err := c.generateSessionID(context.Background(), account.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成sessionID错误"})
		return
	}

	//返回信息
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok",
		"data": &LoginResp{
			SessionID: sessionID,
			UserID:    account.UserID,
			Nickname:  account.Nickname,
		},
	})

	return
}

func (c *CmsApp) generateSessionID(ctx context.Context, userID string) (string, error) {
	// session id 生成
	sessionID := uuid.New().String()

	// session id 存储到数据库
	//key : session_id:{user_id}    val: session_id   8h
	// 以用户设置有效时间
	sessionKey := utils.GetSessionKey(userID)
	err := c.rdb.Set(ctx, sessionKey, sessionID, 8*time.Hour).Err()
	if err != nil {
		fmt.Println("rdb set error: ", err)
		return "", err
	}
	// 以session设置有效时间
	authKey := utils.GetAuthKey(sessionID)
	err = c.rdb.Set(ctx, authKey, time.Now().Unix(), 1*time.Hour).Err()
	if err != nil {
		fmt.Println("rdb set error: ", err)
		return "", err
	}
	return sessionID, nil
}
