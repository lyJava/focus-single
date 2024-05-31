package controller

import (
	"context"
	"errors"
	"focus-single/api/v1"
	"focus-single/internal/model"
	"focus-single/internal/service"
	"focus-single/internal/util"
	"github.com/gogf/gf/v2/frame/g"
	"log"
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
	detailDataOut, err := service.Content().GetDetail(ctx, req.Id)
	if err != nil {
		log.Printf("更新获取数据错误==%+v", err)
		return nil, errors.New("获取数据错误")
	}
	if detailDataOut == nil {
		log.Println("无数据")
		return nil, errors.New("无数据")
	}

	g.Log().Infof(ctx, "更新获取数据==%s", util.ToJsonFormat(detailDataOut, false))

	//getUser := service.Session().GetUser(ctx)
	//g.Log().Infof(ctx, "更新获取用户数据==%s", util.ToJsonFormat(getUser, true))

	service.View().Render(ctx, model.View{
		ContentType: detailDataOut.Content.Type,
		Data:        detailDataOut,
		Title: service.View().GetTitle(ctx, &model.ViewGetTitleInput{
			ContentType: detailDataOut.Content.Type,
			CategoryId:  detailDataOut.Content.CategoryId,
			CurrentName: detailDataOut.Content.Title,
		}),
		BreadCrumb: service.View().GetBreadCrumb(ctx, &model.ViewGetBreadCrumbInput{
			ContentId:   detailDataOut.Content.Id,
			ContentType: detailDataOut.Content.Type,
			CategoryId:  detailDataOut.Content.CategoryId,
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
