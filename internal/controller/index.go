package controller

import (
	"context"
	"log"

	"focus-single/api/v1"
	"focus-single/internal/model"
	"focus-single/internal/service"
)

// Index 首页接口
var Index = cIndex{}

type cIndex struct{}

func (a *cIndex) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	log.Printf("首页类型:%s", req.Type)
	if getListRes, err := service.Content().GetList(ctx, model.ContentGetListInput{
		Type:       req.Type,
		CategoryId: req.CategoryId,
		Page:       req.Page,
		Size:       req.Size,
		Sort:       req.Sort,
	}); err != nil {
		return nil, err
	} else {
		service.View().Render(ctx, model.View{
			ContentType: req.Type,
			Data:        getListRes,
			Title:       "首页",
		})
	}
	return
}
