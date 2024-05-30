package reply

import (
	"context"
	"focus-single/internal/dao"
	"focus-single/internal/model"
	"focus-single/internal/model/entity"
	"focus-single/internal/service"
	"focus-single/internal/util"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"
)

type sReply struct{}

func init() {
	service.RegisterReply(New())
}

func New() *sReply {
	return &sReply{}
}

// Create 创建回复
func (s *sReply) Create(ctx context.Context, in model.ReplyCreateInput) error {
	return dao.Reply.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 覆盖用户ID
		in.UserId = service.BizCtx().Get(ctx).User.Id
		_, err := dao.Reply.Ctx(ctx).TX(tx).Data(in).Insert()
		if err == nil {
			err = service.Content().AddReplyCount(ctx, in.TargetId, 1)
		}
		if err != nil {
			g.Log().Errorf(ctx, "创建回复错误===%+v", err)
		}
		return err
	})
}

// Delete 删除回复(硬删除)
func (s *sReply) Delete(ctx context.Context, id uint) error {
	return dao.Reply.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var reply *entity.Reply
		err := dao.Reply.Ctx(ctx).WherePri(id).Scan(&reply)
		if err != nil {
			g.Log().Errorf(ctx, "查询回复错误===%+v", err)
			return err
		}

		replyContent := reply.Content
		var replyContentSrcList []string
		if replyContent != "" {
			replyContentSrcList, err = util.FindImgSrc(replyContent)
			if err != nil {
				g.Log().Errorf(ctx, "获取图片src错误===%+v", err)
			}
		}

		// 删除回复记录
		_, err = dao.Reply.Ctx(ctx).TX(tx).
			Where(dao.Reply.Columns().Id, id).
			Where(dao.Reply.Columns().UserId, service.BizCtx().Get(ctx).User.Id).
			Delete()
		if err == nil {
			// 回复统计-1
			err = service.Content().AddReplyCount(ctx, reply.TargetId, -1)
			if err != nil {
				g.Log().Errorf(ctx, "回复统计操作错误===%+v", err)
				return err
			}
			// 判断回复是否采纳
			var content *entity.Content
			err = dao.Content.Ctx(ctx).TX(tx).Where("id", reply.TargetId).Scan(&content)
			if err == nil && content != nil && content.AdoptedReplyId == id {
				err = service.Content().UnacceptedReply(ctx, reply.TargetId)
			}
		}
		if err != nil {
			g.Log().Errorf(ctx, "删除回复错误===%+v", err)
			return err
		}

		// 删除回复图片
		g.Log().Infof(ctx, "获取图片src切片:%v", replyContentSrcList)
		err = util.DeleteFile(replyContentSrcList)
		if err != nil {
			g.Log().Errorf(ctx, "删除图片出错===%+v", err)
		}

		return nil
	})
}

// DeleteByUserContentId 删除回复(硬删除)
func (s *sReply) DeleteByUserContentId(ctx context.Context, userId, contentId uint) error {
	return dao.Reply.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 删除内容对应的回复
		_, err := dao.Reply.Ctx(ctx).TX(tx).Where(g.Map{
			dao.Reply.Columns().TargetId: contentId,
			dao.Reply.Columns().UserId:   userId,
		}).Delete()
		if err != nil {
			g.Log().Errorf(ctx, "删除回复(硬删除)错误===%+v", err)
		}
		return err
	})
}

// GetList 获取回复列表
func (s *sReply) GetList(ctx context.Context, in model.ReplyGetListInput) (out *model.ReplyGetListOutput, err error) {
	out = &model.ReplyGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}

	m := dao.Reply.Ctx(ctx).Fields(model.ReplyListItem{})
	if in.TargetType != "" {
		m = m.Where(dao.Reply.Columns().TargetType, in.TargetType)
	}
	if in.TargetId > 0 {
		m = m.Where(dao.Reply.Columns().TargetId, in.TargetId)
	}
	if in.UserId > 0 {
		m = m.Where(dao.Reply.Columns().UserId, in.UserId)
	}

	if err = m.Page(in.Page, in.Size).OrderDesc(dao.Content.Columns().Id).ScanList(&out.List, "Reply"); err != nil {
		g.Log().Errorf(ctx, "获取回复列表错误===%+v", err)
		return nil, err
	}
	if len(out.List) == 0 {
		g.Log().Info(ctx, "回复列表为空")
		return nil, nil
	}

	userIdList := gutil.ListItemValuesUnique(out.List, "Reply", "UserId")
	targetIdList := gutil.ListItemValuesUnique(out.List, "Reply", "TargetId")

	g.Log().Printf(ctx, "最开始的targetIdList==%v,userIdList==%v", targetIdList, userIdList)

	// 用户信息
	if err = dao.User.Ctx(ctx).
		Fields(model.ReplyListUserItem{}).
		WhereIn(dao.User.Columns().Id, userIdList).
		ScanList(&out.List, "User", "Reply", "id:UserId"); err != nil {
		g.Log().Errorf(ctx, "获取用户信息错误===%+v", err)
		return nil, err
	}

	// 内容信息
	if err = dao.Content.Ctx(ctx).
		Fields(dao.Content.Columns().Id, dao.Content.Columns().Title, dao.Content.Columns().CategoryId).
		WhereIn(dao.Content.Columns().Id, targetIdList).
		ScanList(&out.List, "Content", "Reply", "id:TargetId"); err != nil {
		g.Log().Errorf(ctx, "获取回复内容错误===%+v", err)
		return nil, err
	}

	// 现场才能正常获取分类ID切片
	categoryIdList := gutil.ListItemValuesUnique(out.List, "Content", "CategoryId")
	g.Log().Printf(ctx, "提取到的categoryIdList==%v", categoryIdList)

	// 类别信息
	if err = dao.Category.Ctx(ctx).
		Fields(model.ContentListCategoryItem{}).
		WhereIn(dao.Category.Columns().Id, categoryIdList).
		ScanList(&out.List, "Category", "Content", "id:CategoryId"); err != nil {
		g.Log().Errorf(ctx, "获取分类信息错误===%+v", err)
		return nil, err
	}
	return out, nil
}
