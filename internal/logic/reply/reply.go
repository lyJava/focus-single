package reply

import (
	"context"
	"encoding/json"
	"focus-single/internal/dao"
	"focus-single/internal/model"
	"focus-single/internal/model/entity"
	"focus-single/internal/service"
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
		// 删除回复记录
		_, err = dao.Reply.Ctx(ctx).TX(tx).Where(g.Map{
			dao.Reply.Columns().Id:     id,
			dao.Reply.Columns().UserId: service.BizCtx().Get(ctx).User.Id,
		}).Delete()
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
		g.Log().Errorf(ctx, "删除回复错误===%+v", err)
		return err
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

	orderDesc := m.Page(in.Page, in.Size).OrderDesc(dao.Content.Columns().Id)
	/*result, _ := orderDesc.All()

	resultList := result.List()
	outMarshal2, _ := json.MarshalIndent(resultList, "", "    ")
	g.Log().Printf(ctx, "最开始的resultList==%s", string(outMarshal2))

	mapList := result.MapKeyUint("id")
	outMarshal1, _ := json.MarshalIndent(mapList, "", "    ")
	g.Log().Printf(ctx, "最开始的mapList==%s", string(outMarshal1))

	targetIdList := gutil.ListItemValuesUnique(resultList, "target_id")
	g.Log().Printf(ctx, "最开始的targetIdList==%v", targetIdList)

	userIdList := gutil.ListItemValuesUnique(resultList, "user_id")
	g.Log().Printf(ctx, "最开始的userIdList==%v", userIdList)*/

	err = orderDesc.ScanList(&out.List, "Reply")
	if err != nil {
		return nil, err
	}
	if len(out.List) == 0 {
		return nil, nil
	}

	outMarshal2, _ := json.MarshalIndent(&out.List, "", "    ")
	g.Log().Printf(ctx, "最开始的out-List==%s", string(outMarshal2))

	userIdList := gutil.ListItemValuesUnique(out.List, "Reply", "UserId")
	targetIdList := gutil.ListItemValuesUnique(out.List, "Reply", "TargetId")
	categoryIdList := gutil.ListItemValuesUnique(out.List, "Content", "CategoryId")

	g.Log().Printf(ctx, "最开始的targetIdList2==%v,userIdList2==%v,categoryIdList==%v", targetIdList, userIdList, categoryIdList)

	// 用户信息
	if err = m.ScanList(&out.List, "Reply"); err != nil {
		return nil, err
	}
	err = dao.User.Ctx(ctx).
		Fields(model.ReplyListUserItem{}).
		WhereIn(dao.User.Columns().Id, userIdList).
		ScanList(&out.List, "User", "Reply", "id:UserId")
	if err != nil {
		return nil, err
	}

	// 内容信息
	err = dao.Content.Ctx(ctx).
		Fields(dao.Content.Columns().Id, dao.Content.Columns().Title, dao.Content.Columns().CategoryId).
		WhereIn(dao.Content.Columns().Id, targetIdList).
		ScanList(&out.List, "Content", "Reply", "id:TargetId")
	if err != nil {
		return nil, err
	}

	if len(categoryIdList) == 0 {
		categoryIdList = gutil.ListItemValuesUnique(&out.List, "Content", "CategoryId")
	}

	// 类别信息
	err = dao.Category.Ctx(ctx).
		Fields(model.ContentListCategoryItem{}).
		WhereIn(dao.Category.Columns().Id, categoryIdList).
		ScanList(&out.List, "Category", "Content", "id:CategoryId")
	if err != nil {
		return nil, err
	}

	/*outList, _ := json.Marshal(&out.List)
	g.Log().Printf(ctx, "最开始的outList==%s", string(outList))

	// 用户信息
	userMap := make(map[uint]*model.ReplyListUserItem)
	var users []*model.ReplyListUserItem
	if len(userIdList) > 0 {
		err = dao.User.Ctx(ctx).
			Fields(model.ReplyListUserItem{}).
			//Where(dao.User.Columns().Id+" IN(?)", userIds).
			WhereIn(dao.User.Columns().Id, userIdList).
			Scan(&users)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			userMap[user.Id] = user
		}
	}

	// 内容信息
	contentMap := make(map[uint]*model.ContentListItem)
	var contents []*model.ContentListItem
	if len(targetIdList) > 0 {
		//gp := []string{"id", "type", "user_id", "title", "category_id", "content"}
		err = dao.Content.Ctx(ctx).
			Fields(
				dao.Content.Columns().Id,
				dao.Content.Columns().Type,
				dao.Content.Columns().CategoryId,
				dao.Content.Columns().UserId,
				dao.Content.Columns().Title,
				dao.Content.Columns().Content,
				dao.Content.Columns().Sort,
			).
			//Where(dao.Content.Columns().Id+" IN(?)", targetIdList).
			//.Group(gp...).
			WhereIn(dao.Content.Columns().Id, targetIdList).
			Scan(&contents)
		if err != nil {
			return nil, err
		}
		for _, content := range contents {
			contentMap[content.Id] = content
		}
	}

	// 分类信息
	categoryIds := gutil.ListItemValuesUnique(contents, "CategoryId")
	categoryMap := make(map[uint]*model.ContentListCategoryItem)
	if len(categoryIds) > 0 {
		var categories []*model.ContentListCategoryItem
		err = dao.Category.Ctx(ctx).
			Fields(
				dao.Category.Columns().Id,
				dao.Category.Columns().Name,
				dao.Category.Columns().ContentType,
				dao.Category.Columns().Thumb,
			).
			WhereIn(dao.Category.Columns().Id, categoryIds).
			Scan(&categories)
		if err != nil {
			return nil, err
		}

		for _, category := range categories {
			categoryMap[category.Id] = category
		}
	}

	userMarshal, _ := json.Marshal(userMap)
	contentMarshal, _ := json.Marshal(contentMap)
	categoryMarshal, _ := json.Marshal(categoryMap)

	g.Log().Printf(ctx, "用户信息Map===%s", string(userMarshal))
	g.Log().Printf(ctx, "回复内容Map===%s", string(contentMarshal))
	g.Log().Printf(ctx, "分类信息Map===%s", string(categoryMarshal))

	for _, item := range out.List {
		itemMarshal, _ := json.MarshalIndent(item, "", "    ")
		g.Log().Printf(ctx, "集合选项item===%s", itemMarshal)
		userId := item.Reply.UserId
		targetId := item.Reply.TargetId
		contentId := item.Reply.Id
		g.Log().Printf(ctx, "用户ID===%d", userId)
		g.Log().Printf(ctx, "内容ID===%d", targetId)
		g.Log().Printf(ctx, "回复ID===%d", contentId)
		item.User = userMap[userId]
		item.Content = contentMap[targetId]
	}

	marshal3, _ := json.MarshalIndent(&out, "", "    ")
	g.Log().Printf(ctx, "out===%s", string(marshal3))*/
	return out, nil
}

// 获取回复列表
func (s *sReply) GetList_old(ctx context.Context, in model.ReplyGetListInput) (out *model.ReplyGetListOutput, err error) {
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

	err = m.Page(in.Page, in.Size).OrderDesc(dao.Content.Columns().Id).ScanList(&out.List, "Reply")
	if err != nil {
		return nil, err
	}
	if len(out.List) == 0 {
		return nil, nil
	}
	// User
	if err = m.ScanList(&out.List, "Reply"); err != nil {
		return nil, err
	}
	err = dao.User.Ctx(ctx).
		Fields(model.ReplyListUserItem{}).
		Where(dao.User.Columns().Id, gutil.ListItemValuesUnique(out.List, "Reply", "UserId")).
		ScanList(&out.List, "User", "Reply", "id:UserId")
	if err != nil {
		return nil, err
	}

	// Content
	err = dao.Content.Ctx(ctx).
		Fields(dao.Content.Columns().Id, dao.Content.Columns().Title, dao.Content.Columns().CategoryId).
		Where(dao.Content.Columns().Id, gutil.ListItemValuesUnique(out.List, "Reply", "TargetId")).
		ScanList(&out.List, "Content", "Reply", "id:TargetId")
	if err != nil {
		return nil, err
	}

	// Category
	err = dao.Category.Ctx(ctx).
		Fields(model.ContentListCategoryItem{}).
		Where(dao.Category.Columns().Id, gutil.ListItemValuesUnique(out.List, "Content", "CategoryId")).
		ScanList(&out.List, "Category", "Content", "id:CategoryId")
	if err != nil {
		return nil, err
	}

	return out, nil
}
