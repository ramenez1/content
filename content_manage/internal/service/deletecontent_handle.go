package service

import (
	"content_manage/api/operate"
	"context"
)

func (a *AppService) DeleteContent(ctx context.Context, req *operate.DeleteContentReq) (*operate.DeleteContentRsp, error) {
	uc := a.uc
	err := uc.DeleteContent(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	return &operate.DeleteContentRsp{}, nil
}
