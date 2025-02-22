package service

import (
	"ContentSystem/internal/dao"
	"ContentSystem/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type RegisterReq struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type RegisterResp struct {
	Message string `json:"message" binding:"required"`
}

func (c *CmsApp) Register(ctx *gin.Context) {
	var req RegisterReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//密码加密
	hashedPassword, err := encryptPassword(req.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{})
	}

	fmt.Printf("register req = %+v\n , hashedPassword = [%v]\n", req, hashedPassword)
	//账号校验
	accountDao := dao.NewAccountDao(c.db)
	isExist, err := accountDao.IsExist(req.UserID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error isExist": err.Error()})
		return
	}
	if isExist {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}
	//信息持久化
	nowTime := time.Now()
	if err := accountDao.CreateAccount(model.Account{
		UserID:     req.UserID,
		Password:   hashedPassword,
		Nickname:   req.Nickname,
		Creattime:  nowTime,
		Updatetime: nowTime,
	}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("register req = %+v\n", req)

	//返回信息
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "ok ",
		"data": &RegisterResp{
			Message: fmt.Sprintf("用户%s注册成功", req.UserID),
		},
	})
}

func encryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("bcrypt generate from password error = %v", err)
		return "", err
	}
	return string(hashedPassword), err
}
