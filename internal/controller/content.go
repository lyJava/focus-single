package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"log"

	"focus-single/api/v1"
	"focus-single/internal/model"
	"focus-single/internal/service"
)

// Content 内容管理
var Content = cContent{}

type cContent struct{}

func (a *cContent) ShowCreate(ctx context.Context, req *v1.ContentShowCreateReq) (res *v1.ContentShowCreateRes, err error) {
	service.View().Render(ctx, model.View{
		ContentType: req.Type,
	})
	return
}

func (a *cContent) Create(ctx context.Context, req *v1.ContentCreateReq) (res *v1.ContentCreateRes, err error) {
	out, err := service.Content().Create(ctx, model.ContentCreateInput{
		ContentCreateUpdateBase: model.ContentCreateUpdateBase{
			Type:       req.Type,
			CategoryId: req.CategoryId,
			Title:      req.Title,
			Content:    req.Content,
			Brief:      req.Brief,
			Thumb:      req.Thumb,
			Tags:       req.Tags,
			Referer:    req.Referer,
		},
		UserId: service.Session().GetUser(ctx).Id,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ContentCreateRes{ContentId: out.ContentId}, nil
}

func (a *cContent) ShowUpdate(ctx context.Context, req *v1.ContentShowUpdateReq) (res *v1.ContentShowUpdateRes, err error) {
	log.Printf("更新获取数据参数==%d", req.Id)
	getDetailRes, err := service.Content().GetDetail(ctx, req.Id)
	if err != nil {
		log.Printf("更新获取数据错误==%+v", err)
		return nil, errors.New("获取数据错误")
	}
	if getDetailRes == nil {
		log.Println("无数据")
		return nil, errors.New("无数据")
	}

	var list []interface{}
	list = append(list, getDetailRes)
	marshal, err := json.Marshal(&list)
	g.Log().Infof(ctx, "更新获取数据==%s", string(marshal))

	getUser := service.Session().GetUser(ctx)
	g.Log().Infof(ctx, "从会话中获取用户数据==%v", getUser)

	getTitle := service.View().GetTitle(ctx, &model.ViewGetTitleInput{
		ContentType: getDetailRes.Content.Type,
		CategoryId:  getDetailRes.Content.CategoryId,
		CurrentName: getDetailRes.Content.Title,
	})

	g.Log().Infof(ctx, "从会话中获取getTitle数据==%v", getTitle)

	service.View().Render(ctx, model.View{
		ContentType: getDetailRes.Content.Type,
		Data: map[string]interface{}{
			"List": list,
		},
		Title: getTitle,
		BreadCrumb: service.View().GetBreadCrumb(ctx, &model.ViewGetBreadCrumbInput{
			ContentId:   getDetailRes.Content.Id,
			ContentType: getDetailRes.Content.Type,
			CategoryId:  getDetailRes.Content.CategoryId,
		}),
	})
	return
}

func (a *cContent) Update(ctx context.Context, req *v1.ContentUpdateReq) (res *v1.ContentUpdateRes, err error) {
	err = service.Content().Update(ctx, model.ContentUpdateInput{
		Id: req.Id,
		ContentCreateUpdateBase: model.ContentCreateUpdateBase{
			Type:       req.Type,
			CategoryId: req.CategoryId,
			Title:      req.Title,
			Content:    req.Content,
			Brief:      req.Brief,
			Thumb:      req.Thumb,
			Tags:       req.Tags,
			Referer:    req.Referer,
		},
	})
	if err != nil {
		log.Printf("更新错误==%+v", err)
	}
	return
}

func (a *cContent) Delete(ctx context.Context, req *v1.ContentDeleteReq) (res *v1.ContentDeleteRes, err error) {
	err = service.Content().Delete(ctx, req.Id)
	if err != nil {
		log.Printf("删除出错误==%+v", err)
	}
	return
}

// AdoptReply 采纳回复
func (a *cContent) AdoptReply(ctx context.Context, req *v1.ReplayAdoptReq) (res *v1.ReplayAdoptRes, err error) {
	if err = service.Content().AdoptReply(ctx, req.Id, req.ReplyId); err != nil {
		log.Printf("采纳回复出错==%+v", err)
	}
	log.Println("采纳回复成功")
	return
}
