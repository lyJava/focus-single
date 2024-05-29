package controller

import (
	"context"
	"log"

	"focus-single/api/v1"
	"focus-single/internal/service"
)

// Interact 赞踩控制器
var Interact = cInteract{}

type cInteract struct{}

// Zan 赞
func (a *cInteract) Zan(ctx context.Context, req *v1.InteractZanReq) (res *v1.InteractZanRes, err error) {
	err = service.Interact().Zan(ctx, req.Type, req.Id)
	if err == nil {
		log.Println("点赞成功了")
	}
	return
}

// CancelZan 取消赞
func (a *cInteract) CancelZan(ctx context.Context, req *v1.InteractCancelZanReq) (res *v1.InteractCancelZanRes, err error) {
	err = service.Interact().CancelZan(ctx, req.Type, req.Id)
	if err == nil {
		log.Println("取消赞成功了")
	}
	return
}

// Cai 踩
func (a *cInteract) Cai(ctx context.Context, req *v1.InteractCaiReq) (res *v1.InteractCaiRes, err error) {
	err = service.Interact().Cai(ctx, req.Type, req.Id)
	if err == nil {
		log.Println("踩成功了")
	}
	return
}

// CancelCai 取消踩
func (a *cInteract) CancelCai(ctx context.Context, req *v1.InteractCancelCaiReq) (res *v1.InteractCancelCaiRes, err error) {
	err = service.Interact().CancelCai(ctx, req.Type, req.Id)
	if err == nil {
		log.Println("取消踩成功了")
	}
	return
}
