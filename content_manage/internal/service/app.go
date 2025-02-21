package service

import (
	"content_manage/api/operate"
	"content_manage/internal/biz"
)

type AppService struct {
	//继承状态码内容
	operate.UnimplementedAppServer

	uc *biz.ContentUsecase
}

// NewAppService new app service
func NewAppService(uc *biz.ContentUsecase) *AppService {
	return &AppService{uc: uc}
}
