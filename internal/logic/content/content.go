package content

import (
	"context"
	"fmt"
	"focus-single/internal/util"
	"log"

	"focus-single/internal/service"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/ghtml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"

	"focus-single/internal/consts"
	"focus-single/internal/dao"
	"focus-single/internal/model"
	"focus-single/internal/model/entity"

	"github.com/spf13/cast"
)

type sContent struct{}

func init() {
	service.RegisterContent(New())
}

func New() *sContent {
	return &sContent{}
}

// GetList 查询内容列表
func (s *sContent) GetList(ctx context.Context, in model.ContentGetListInput) (out *model.ContentGetListOutput, err error) {
	//log.Println("首页查询")
	m := dao.Content.Ctx(ctx)
	out = &model.ContentGetListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	// 默认查询topic
	if in.Type != "" {
		m = m.Where(dao.Content.Columns().Type, in.Type)
	} else {
		m = m.Where(dao.Content.Columns().Type, consts.ContentTypeTopic)
	}
	// 栏目检索
	if in.CategoryId > 0 {
		idArray, err := service.Category().GetSubIdList(ctx, in.CategoryId)
		if err != nil {
			return out, err
		}
		m = m.Where(dao.Content.Columns().CategoryId, idArray)
	}
	// 管理员可以查看所有文章
	if in.UserId > 0 && !service.User().IsAdmin(ctx, in.UserId) {
		m = m.Where(dao.Content.Columns().UserId, in.UserId)
	}
	// 分配查询
	listModel := m.Page(in.Page, in.Size)
	// 排序方式
	switch in.Sort {
	case consts.ContentSortHot:
		listModel = listModel.OrderDesc(dao.Content.Columns().ViewCount)

	case consts.ContentSortActive:
		listModel = listModel.OrderDesc(dao.Content.Columns().UpdatedAt)

	default:
		listModel = listModel.OrderDesc(dao.Content.Columns().Id)
	}
	// 执行查询
	var list []*entity.Content
	if err := listModel.Scan(&list); err != nil {
		return out, err
	}
	// 没有数据
	if len(list) == 0 {
		return out, nil
	}
	out.Total, err = m.Count()
	if err != nil {
		return out, err
	}
	// Content
	if err := listModel.ScanList(&out.List, "Content"); err != nil {
		return out, err
	}
	// Category
	err = dao.Category.Ctx(ctx).
		Fields(model.ContentListCategoryItem{}).
		Where(dao.Category.Columns().Id, gutil.ListItemValuesUnique(out.List, "Content", "CategoryId")).
		ScanList(&out.List, "Category", "Content", "id:CategoryId")
	if err != nil {
		return out, err
	}
	// User
	err = dao.User.Ctx(ctx).
		Fields(model.ContentListUserItem{}).
		Where(dao.User.Columns().Id, gutil.ListItemValuesUnique(out.List, "Content", "UserId")).
		ScanList(&out.List, "User", "Content", "id:UserId")
	if err != nil {
		return out, err
	}
	//marshal, _ := json.Marshal(&out.List)
	//g.Log().Infof(ctx, "首页查询list===%s", string(marshal))
	return
}

// Search 搜索内容列表
func (s *sContent) Search(ctx context.Context, in model.ContentSearchInput) (out *model.ContentSearchOutput, err error) {
	var (
		m           = dao.Content.Ctx(ctx)
		likePattern = `%` + in.Key + `%`
	)
	out = &model.ContentSearchOutput{
		Page: in.Page,
		Size: in.Size,
	}
	m = m.WhereLike(dao.Content.Columns().Content, likePattern).WhereOrLike(dao.Content.Columns().Title, likePattern)
	// 内容类型
	if in.Type != "" {
		m = m.Where(dao.Content.Columns().Type, in.Type)
	}
	// 栏目检索
	if in.CategoryId > 0 {
		idArray, err := service.Category().GetSubIdList(ctx, in.CategoryId)
		if err != nil {
			return nil, err
		}
		m = m.Where(dao.Content.Columns().CategoryId, idArray)
	}
	listModel := m.Page(in.Page, in.Size)
	switch in.Sort {
	case consts.ContentSortHot:
		listModel = listModel.OrderDesc(dao.Content.Columns().ViewCount)

	case consts.ContentSortActive:
		listModel = listModel.OrderDesc(dao.Content.Columns().UpdatedAt)

	// case model.ContentSortScore:
	//	listModel = listModel.OrderDesc("score")

	default:
		listModel = listModel.OrderDesc(dao.Content.Columns().Id)
	}
	all, err := listModel.All()
	if err != nil {
		return nil, err
	}
	// 没有数据
	if all.IsEmpty() {
		return out, nil
	}
	out.Total, err = m.Count()
	if err != nil {
		return nil, err
	}
	// 搜索统计
	statsModel := m.Fields(dao.Content.Columns().Type, "count(*) total").Group(dao.Content.Columns().Type)
	statsAll, err := statsModel.All()
	if err != nil {
		return nil, err
	}
	out.Stats = make(map[string]int)
	for _, v := range statsAll {
		out.Stats[v["type"].String()] = v["total"].Int()
	}
	// Content
	if err = all.ScanList(&out.List, "Content"); err != nil {
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
	// User
	err = dao.User.Ctx(ctx).
		Fields(model.ContentListUserItem{}).
		Where(dao.User.Columns().Id, gutil.ListItemValuesUnique(out.List, "Content", "UserId")).
		ScanList(&out.List, "User", "Content", "id:UserId")
	if err != nil {
		return nil, err
	}

	return out, nil
}

// GetDetail 查询详情
func (s *sContent) GetDetail(ctx context.Context, id uint) (out *model.ContentGetDetailOutput, err error) {
	out = &model.ContentGetDetailOutput{}
	if err := dao.Content.Ctx(ctx).
		Where(dao.Content.Columns().Id, id).
		Scan(&out.Content); err != nil {
		return nil, err
	}
	// 没有数据
	if out.Content == nil {
		return nil, nil
	}
	if err = dao.User.Ctx(ctx).Where(dao.Content.Columns().Id, out.Content.UserId).Scan(&out.User); err != nil {
		return nil, err
	}

	// 查询分类数据
	if err = dao.Category.Ctx(ctx).
		Fields(dao.Category.Columns().Id, dao.Category.Columns().Name, dao.Category.Columns().ContentType).
		Where(dao.Category.Columns().Id, out.Content.CategoryId).Scan(&out.Category); err != nil {
		return nil, err
	}

	//marshal, _ := json.Marshal(&out)
	//g.Log().Infof(ctx, "返回的详情:%s", string(marshal))
	return out, nil
}

// GetDetail2 查询详情
func (s *sContent) GetDetail2(ctx context.Context, id, categoryId uint) (out *model.ContentGetDetailOutput, err error) {
	out = &model.ContentGetDetailOutput{}
	if err = dao.Content.Ctx(ctx).
		Where(dao.Content.Columns().Id, id).
		Where(dao.Content.Columns().CategoryId, categoryId).
		Scan(&out.Content); err != nil {
		return nil, err
	}
	// 没有数据
	if out.Content == nil {
		return nil, nil
	}
	if err = dao.User.Ctx(ctx).Where(dao.Content.Columns().Id, out.Content.UserId).Scan(&out.User); err != nil {
		return nil, err
	}

	// 查询分类数据
	if err = dao.Category.Ctx(ctx).
		Fields(dao.Category.Columns().Id, dao.Category.Columns().Name, dao.Category.Columns().ContentType).
		Where(dao.Category.Columns().Id, out.Content.CategoryId).Scan(&out.Category); err != nil {
		return nil, err
	}
	return out, nil
}

// Create 创建内容
//
//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection,SqlDialectInspection,GoConvertStringLiterals
func (s *sContent) Create(ctx context.Context, in model.ContentCreateInput) (out model.ContentCreateOutput, err error) {
	if in.UserId == 0 {
		in.UserId = service.BizCtx().Get(ctx).User.Id
	}
	// 不允许HTML代码
	if err = ghtml.SpecialCharsMapOrStruct(in); err != nil {
		return out, err
	}

	sqlType := dao.Interact.DB().GetConfig().Type
	log.Printf("数据库类型===%s", sqlType)

	var lastRowId int64

	if sqlType == "mysql" {
		err = dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			lastRowId, err = dao.Content.Ctx(ctx).TX(tx).Data(in).InsertAndGetId()
			if err != nil {
				log.Printf("发布内容出错===%+v", err)
				return err
			}
			return nil
		})
		if err != nil {
			log.Printf("发布内容出错===%+v", err)
			return out, err
		}
		return out, nil
	}

	if sqlType == "pgsql" {
		insertSql := `INSERT INTO gf_content(user_id, type, category_id,  title, content, created_at, updated_at) 
    					VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`
		err = g.DB().Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			query, err := g.DB().Ctx(ctx).GetOne(ctx, insertSql, in.UserId, in.Type, in.CategoryId, in.Title, in.Content)
			if err != nil {
				log.Printf("发布内容出错===%+v", err)
				return err
			}
			lastRowId = cast.ToInt64(query.Map()["id"])
			return nil
		})
		if err != nil {
			return model.ContentCreateOutput{}, err
		}
	}

	log.Printf("发布内容返回===%d", lastRowId)

	return model.ContentCreateOutput{ContentId: uint(lastRowId)}, nil
}

// Update 修改
func (s *sContent) Update(ctx context.Context, in model.ContentUpdateInput) error {
	return dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 不允许HTML代码
		if err := ghtml.SpecialCharsMapOrStruct(in); err != nil {
			return err
		}
		_, err := dao.Content.Ctx(ctx).TX(tx).
			Data(in).
			FieldsEx(dao.Content.Columns().Id).
			Where(dao.Content.Columns().Id, in.Id).
			Where(dao.Content.Columns().UserId, service.BizCtx().Get(ctx).User.Id).
			Update()
		return err
	})
}

// Delete 删除
func (s *sContent) Delete(ctx context.Context, id uint) error {
	return dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		user := service.BizCtx().Get(ctx).User
		if user == nil {
			return fmt.Errorf("未登录不允许操作")
		}
		log.Printf("是否是管理员:%t", user.IsAdmin)
		// 管理员直接删除文章和评论
		if user.IsAdmin {
			var contentEntity *entity.Content
			err := dao.Content.Ctx(ctx).TX(tx).Where(dao.Content.Columns().Id, id).Scan(&contentEntity)
			if err == nil && contentEntity != nil {
				contentStr := contentEntity.Content
				log.Printf("对应内容信息htmlStr:%s", contentStr)
				if contentStr != "" {
					contentImgSrcList := util.GetImgSrcFromStr(contentStr)
					log.Printf("内容图片切片===%v", contentImgSrcList)
					err = util.DeleteFile(contentImgSrcList)
					if err != nil {
						return err
					}
				}
			}

			_, err = dao.Content.Ctx(ctx).TX(tx).Where(dao.Content.Columns().Id, id).Delete()
			log.Printf("内容ID:%d", id)
			if err == nil {
				var replyEntity *entity.Reply
				// TargetId为内容表的主键ID
				err = dao.Reply.ReplyDao.Ctx(ctx).TX(tx).Where(dao.Reply.Columns().TargetId, id).Scan(&replyEntity)
				if err == nil && replyEntity != nil {
					contentStr := replyEntity.Content
					log.Printf("对应回复信息htmlStr:%s", contentStr)
					if contentStr != "" {
						var replyImgSrcList []string
						replyImgSrcList, err = util.FindImgSrc(contentStr)
						if err == nil {
							log.Printf("回复内容图片切片===%v", replyImgSrcList)
							err = util.DeleteFile(replyImgSrcList)
							if err != nil {
								return err
							}
						}
					}
				}

				_, err = dao.Reply.Ctx(ctx).TX(tx).Where(dao.Reply.Columns().TargetId, id).Delete()
			}
			return err
		}
		// 删除内容
		_, err := dao.Content.Ctx(ctx).TX(tx).Where(g.Map{
			dao.Content.Columns().Id:     id,
			dao.Content.Columns().UserId: service.BizCtx().Get(ctx).User.Id,
		}).Delete()
		// 删除评论
		if err != nil {
			return err
		}
		err = service.Reply().DeleteByUserContentId(ctx, user.Id, id)
		if err != nil {
			return err
		}

		return nil
	})
}

// AddViewCount 浏览次数增加
func (s *sContent) AddViewCount(ctx context.Context, id uint, count int) error {
	return dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.Content.Ctx(ctx).TX(tx).
			Where(dao.Content.Columns().Id, id).
			Increment(dao.Content.Columns().ViewCount, count)
		if err != nil {
			return err
		}
		return nil
	})
}

// AddReplyCount 增加回复次数
func (s *sContent) AddReplyCount(ctx context.Context, id uint, count int) error {
	return dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.Content.Ctx(ctx).TX(tx).
			Where(dao.Content.Columns().Id, id).
			Increment(dao.Content.Columns().ReplyCount, count)
		if err != nil {
			return err
		}
		log.Println("更新回复次数成功")
		return nil
	})
}

// AdoptReply 采纳回复
func (s *sContent) AdoptReply(ctx context.Context, id uint, replyID uint) error {
	return dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.Content.Ctx(ctx).TX(tx).
			Data(dao.Content.Columns().AdoptedReplyId, replyID).
			Where(dao.Content.Columns().UserId, service.BizCtx().Get(ctx).User.Id).
			Where(dao.Content.Columns().Id, id).
			Update()
		if err != nil {
			return err
		}
		return nil
	})
}

// UnacceptedReply 取消采纳回复
func (s *sContent) UnacceptedReply(ctx context.Context, id uint) error {
	return dao.Content.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.Content.Ctx(ctx).TX(tx).
			Data(dao.Content.Columns().AdoptedReplyId, 0).
			WherePri(dao.Content.Columns().Id, id).
			Update()
		if err != nil {
			return err
		}
		return nil
	})
}
