package service

import (
	"content_manage/api/operate"
	"content_manage/internal/biz"
	"context"
)

func (a *AppService) FindContent(ctx context.Context, req *operate.FindContentReq) (*operate.FindContentRsp, error) {
	findParams := &biz.FindParams{
		ID:       int(req.GetId()),
		Author:   req.GetAuthor(),
		Title:    req.GetTitle(),
		Page:     int(req.GetPage()),
		PageSize: int(req.GetPageSize()),
	}
	uc := a.uc
	results, total, err := uc.FindContent(ctx, findParams)
	if err != nil {
		return nil, err
	}
	var contents []*operate.Content
	for _, r := range results {
		contents = append(contents, &operate.Content{
			Id:             int64(r.ID),
			Title:          r.Title,
			Description:    r.Description,
			Author:         r.Author,
			VideoUrl:       r.VideoURL,
			Thumbnail:      r.Thumbnail,
			Category:       r.Category,
			Duration:       r.Duration.Milliseconds(),
			Resolution:     r.Resolution,
			FileSize:       r.FileSize,
			Format:         r.Format,
			Quality:        int32(r.Quality),
			ApprovalStatus: int32(r.ApprovalStatus),
		})
	}
	rsp := &operate.FindContentRsp{
		Content: contents,
		Total:   int64(total),
	}
	return rsp, nil
}
