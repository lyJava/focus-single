package view

import (
	"context"
	"encoding/json"
	"focus-single/internal/util"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmode"
	"html/template"

	"focus-single/internal/consts"
	"focus-single/internal/model"
	"focus-single/internal/service"
)

// viewBuildIn 视图自定义方法管理对象
type viewBuildIn struct {
	httpRequest *ghttp.Request
}

// DidIZan 我是否赞了这个内容
func (s *viewBuildIn) DidIZan(targetType string, targetId uint) bool {
	b, _ := service.Interact().DidIZan(s.httpRequest.Context(), targetType, targetId)
	return b
}

// DidICai 我是否踩了这个内容
func (s *viewBuildIn) DidICai(targetType string, targetId uint) bool {
	b, _ := service.Interact().DidICai(s.httpRequest.Context(), targetType, targetId)
	return b
}

// TopMenus 获取顶部菜单列表
func (s *viewBuildIn) TopMenus() ([]*model.MenuItem, error) {
	topMenus, err := service.Menu().GetTopMenus(s.httpRequest.Context())
	if err != nil {
		return nil, err
	}
	if len(topMenus) == 0 {
		return nil, nil
	}
	currentUriWithQueryString := s.httpRequest.URL.String()
	// 处理是否选中, URL，包含QueryString
	for _, v := range topMenus {
		if gstr.Equal(v.Url, currentUriWithQueryString) {
			v.Active = true
			return topMenus, nil
		}
	}
	// 处理是否选中, URI
	for _, v := range topMenus {
		if v.Url == "/" {
			continue
		}
		if gstr.HasPrefix(currentUriWithQueryString, v.Url) {
			v.Active = true
			return topMenus, nil
		}
	}
	// 没有选中的菜单，那么自动识别第一层路由，例如：
	// /topic/1 则选中 /topic 菜单。
	array := gstr.SplitAndTrim(s.httpRequest.URL.Path, "/")
	if len(array) > 1 {
		path := "/" + array[0]
		for _, v := range topMenus {
			if gstr.Equal(v.Url, path) {
				v.Active = true
				return topMenus, nil
			}
		}
	}
	// 最后则自动高亮首页(第一个菜单)
	topMenus[0].Active = true
	return topMenus, nil
}

// CategoryTree 获得指定的栏目树形对象，当contentType为空时，表示获取所有的栏目树形对象。
func (s *viewBuildIn) CategoryTree(contentType string) ([]*model.CategoryTreeItem, error) {
	ctx := context.Background()
	g.Log().Infof(ctx, "获得指定的类型===%s", contentType)
	treeItems, err := service.Category().GetTree(ctx, contentType)

	if err != nil {
		g.Log().Errorf(ctx, "获得指定的栏目树形对象出错===%+v", err)
		return nil, err
	}
	marshal, _ := json.Marshal(&treeItems)
	g.Log().Infof(ctx, "获得指定的栏目===%s", marshal)
	return treeItems, nil
}

// IsNew 最新文章，12小时内为最新文章
func (s *viewBuildIn) IsNew(gt *gtime.Time) bool {
	if gt == nil {
		return false
	}
	n := gtime.Now().Timestamp()
	t := gt.Timestamp()

	// var hs int64 = 3600
	hs := int64(3600)

	// 最新文章
	if n-t < hs*12 {
		return true
	}

	return false
}

// GenderFont 根据性别字段内容返回性别的font。
func (s *viewBuildIn) GenderFont(gender int) string {
	switch gender {
	case consts.UserGenderMale:
		return "&#xe651;"
	case consts.UserGenderFemale:
		return "&#xe636;"
	default:
		return "&#xead2;"
	}
}

// Gender 根据性别字段内容返回性别。
func (s *viewBuildIn) Gender(gender int) string {
	/*switch gender {
	case consts.UserGenderMale:
		return "男"
	case consts.UserGenderFemale:
		return "女"
	default:
		return "未知"
	}*/
	return consts.GetGenderByType(gender)
}

// ContentTypeName 根据性别字段内容返回性别。
func (s *viewBuildIn) ContentTypeName(contentType string) string {
	/*switch contentType {
	case consts.ContentTypeArticle:
		return "文章"
	case consts.ContentTypeAsk:
		return "问答"
	case consts.ContentTypeTopic:
		return "主题"
	default:
		return "未知"
	}*/
	return consts.GetContentByType(contentType)
}

// Page 创建分页HTML内容
func (s *viewBuildIn) Page(total, size int) string {
	page := s.httpRequest.GetPage(total, size)
	page.LinkStyle = "page-link"
	page.SpanStyle = "page-link"
	page.PrevPageTag = "«"
	page.NextPageTag = "»"
	content := page.PrevPage() + page.PageBar() + page.NextPage()
	content = gstr.ReplaceByMap(content, map[string]string{
		"<span":  "<li class=\"page-item disabled\"><span",
		"/span>": "/span></li>",
		"<a":     "<li class=\"page-item\"><a",
		"/a>":    "/a></li>",
	})
	return content
}

// UrlPath 获取当前页面的Url Path.
func (s *viewBuildIn) UrlPath() string {
	return s.httpRequest.URL.Path
}

// FormatTime 格式化时间
func (s *viewBuildIn) FormatTime(gt *gtime.Time) string {
	return util.FormatGfTime(gt)
}

// HTMLEscapeStr 转义 HTML 标签
func (s *viewBuildIn) HTMLEscapeStr(htmlStr string) string {
	return template.HTMLEscapeString(htmlStr)
}

// Version 随机数 开发环境时间戳，线上为前端版本号
func (s *viewBuildIn) Version() string {
	var rand string
	if gmode.IsDevelop() {
		rand = gconv.String(gtime.TimestampMilli())
	} else {
		rand = consts.Version
	}
	return rand
}
