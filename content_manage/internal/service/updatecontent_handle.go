package service

import (
	"content_manage/api/operate"
	"content_manage/internal/biz"
	"context"
	"time"
)

func (a *AppService) UpdateContent(ctx context.Context, req *operate.UpdateContentReq) (*operate.UpdateContentRsp, error) {
	content := req.GetContent()
	uc := a.uc
	err := uc.UpdateContent(ctx, &biz.Content{
		ID:             int(content.GetId()),
		Title:          content.GetTitle(),
		Description:    content.GetDescription(),
		Author:         content.GetAuthor(),
		VideoURL:       content.GetVideoUrl(),
		Thumbnail:      content.GetThumbnail(),
		Category:       content.GetCategory(),
		Duration:       time.Duration(content.GetDuration()),
		Resolution:     content.GetResolution(),
		FileSize:       content.GetFileSize(),
		Format:         content.GetFormat(),
		Quality:        int(content.GetQuality()),
		ApprovalStatus: int(content.GetApprovalStatus()),
	})
	if err != nil {
		return nil, err
	}
	return &operate.UpdateContentRsp{}, nil
}
