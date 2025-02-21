package service

import (
	"bytes"
	"content_manage/api/operate"
	"content_manage/internal/biz"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *AppService) CreateContent(ctx context.Context, req *operate.CreateContentReq) (*operate.CreateContentRsp, error) {
	content := req.GetContent()
	uc := a.uc
	contentID := uuid.New().String()
	_, err := uc.CreateContent(ctx, &biz.Content{
		ContentID:      contentID,
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
	err = a.ExecFlow(contentID)
	if err != nil {
		return nil, err
	}
	return &operate.CreateContentRsp{}, nil
}

func (a *AppService) ExecFlow(contentID string) error {
	url := "http://localhost:7788/content-flow"
	method := "GET"
	payload := map[string]interface{}{
		"content_id": contentID}
	data, _ := json.Marshal(payload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("session_id", "aaaa")
	req.Header.Add("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
