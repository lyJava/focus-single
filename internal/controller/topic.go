package controller

import (
	"context"
	"focus-single/api/v1"
	"focus-single/internal/consts"
	"focus-single/internal/model"
	"focus-single/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"log"
)

// Topic 主题管理
var Topic = cTopic{}

type cTopic struct{}

func (a *cTopic) Index(ctx context.Context, req *v1.TopicIndexReq) (res *v1.TopicIndexRes, err error) {
	req.Type = consts.ContentTypeTopic
	out, err := service.Content().GetList(ctx, model.ContentGetListInput{
		Type:       req.Type,
		CategoryId: req.CategoryId,
		Page:       req.Page,
		Size:       req.Size,
		Sort:       req.Sort,
	})
	if err != nil {
		return nil, err
	}
	title := service.View().GetTitle(ctx, &model.ViewGetTitleInput{
		ContentType: req.Type,
		CategoryId:  req.CategoryId,
	})
	service.View().Render(ctx, model.View{
		ContentType: req.Type,
		Data:        out,
		Title:       title,
	})
	return
}

func (a *cTopic) Detail(ctx context.Context, req *v1.TopicDetailReq) (res *v1.TopicDetailRes, err error) {
	log.Printf("当前主题ID:%d", req.Id)
	out, err := service.Content().GetDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if out == nil {
		service.View().Render404(ctx)
		return
	}

	err = service.Content().AddViewCount(ctx, req.Id, 0)
	service.View().Render(ctx, model.View{
		ContentType: consts.ContentTypeTopic,
		Data:        out,
		Title: service.View().GetTitle(ctx, &model.ViewGetTitleInput{
			ContentType: out.Content.Type,
			CategoryId:  out.Content.CategoryId,
			CurrentName: out.Content.Title,
		}),
		BreadCrumb: service.View().GetBreadCrumb(ctx, &model.ViewGetBreadCrumbInput{
			ContentId:   out.Content.Id,
			ContentType: out.Content.Type,
			CategoryId:  out.Content.CategoryId,
		}),
	})
	return
}

// Delete 删除主题内容(包括了主题内容图片)，那么该内容下的回复也会被删除(包括了回复内容图片)
func (a *cTopic) Delete(ctx context.Context, req *v1.TopicDeleteReq) (res *v1.TopicDetailRes, err error) {
	err = service.Content().Delete(ctx, req.Id)
	if err != nil {
		g.Log().Errorf(ctx, "内容删除错误===%+v", err)
	}
	log.Println("内容删除成功")
	return
}
