package consts

const (
	ContentListDefaultSize = 10
	ContentListMaxSize     = 50
	ContentSortDefault     = 0 // 排序：按照创建时间
	ContentSortActive      = 1 // 排序：按照更新时间
	ContentSortHot         = 2 // 排序：按照浏览量
	ContentSortScore       = 3 // 排序：按照搜索结果关联性
	ContentTypeArticle     = "article"
	ContentTypeAsk         = "ask"
	ContentTypeTopic       = "topic"
)

func GetContentByType(typeStr string) string {
	dataMap := map[string]string{
		ContentTypeArticle: "文章",
		ContentTypeTopic:   "主题",
		ContentTypeAsk:     "问答",
	}
	if content, exists := dataMap[typeStr]; exists {
		return content
	}
	return "未知"
}
