package util

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/net/html"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FormatGfTime 格式化时间
func FormatGfTime(gt *gtime.Time) string {
	if gt == nil {
		return ""
	}

	const (
		yearInSeconds   = 31536000
		dayInSeconds    = 86400
		hourInSeconds   = 3600
		minuteInSeconds = 60
		secondInSeconds = 1
	)

	diff := gtime.Now().Timestamp() - gt.Timestamp()

	switch {
	case diff > yearInSeconds:
		return fmt.Sprintf("%d年前", int(diff/yearInSeconds))
	case diff > dayInSeconds:
		return fmt.Sprintf("%d天前", int(diff/dayInSeconds))
	case diff > hourInSeconds:
		return fmt.Sprintf("%d小时前", int(diff/hourInSeconds))
	case diff > minuteInSeconds:
		return fmt.Sprintf("%d分钟前", int(diff/minuteInSeconds))
	case diff > secondInSeconds:
		return fmt.Sprintf("%d秒前", int(diff/secondInSeconds))
	default:
		return "刚刚"
	}
}

// FindImgSrc 从给定的 HTML 字符串中提取所有 img 标签的 src 属性值
func FindImgSrc(htmlStr string) ([]string, error) {
	if htmlStr == "" {
		return nil, nil
	}

	var srcList []string

	// 解析 HTML 字符串
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, fmt.Errorf("解析html错误: %w", err)
	}

	// 定义递归函数遍历节点树
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "img" {
			for _, attr := range node.Attr {
				if attr.Key == "src" {
					srcList = append(srcList, attr.Val)
					break
				}
			}
		}
		// 递归遍历子节点
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
	}

	// 开始遍历
	traverse(doc)

	return srcList, nil
}

// DeleteFile 删除指定路径中的所有文件
func DeleteFile(paths []string) error {

	if len(paths) == 0 {
		return nil
	}

	ctx := context.Background()

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	g.Log().Infof(ctx, "当前目录:%s", currentDir)

	// 从配置文件中获取文件上传路径
	configUploadPath := g.Cfg().MustGet(ctx, "upload.path")
	pathPrefix := strings.Split(configUploadPath.String(), "/")
	if len(pathPrefix) == 0 {
		return fmt.Errorf("文件上传路径必须统一配置")
	}
	filePathCommon := filepath.Join(currentDir, pathPrefix[0])

	var errs []error
	for _, path := range paths {
		g.Log().Infof(ctx, "尝试删除图片相对路径:%s", path)

		fullPath := filepath.Join(filePathCommon, path)
		g.Log().Infof(ctx, "删除文件绝对路径:%s", fullPath)
		// 检查文件权限
		err = CheckFileAndPermissions(fullPath)
		if err != nil {
			errs = append(errs, fmt.Errorf("file check failed %s: %+v", fullPath, err))
			continue
		}
		err = os.RemoveAll(fullPath)
		if err != nil {
			g.Log().Errorf(ctx, "图片:%s删除错误===%+v", fullPath, err)
			errs = append(errs, fmt.Errorf("failed to delete %s: %w", fullPath, err))
		} else {
			g.Log().Infof(ctx, "图片:%s 删除成功", fullPath)
		}
	}

	if len(errs) > 0 {
		// 汇总所有错误
		var errMsg string
		for _, e := range errs {
			errMsg += e.Error() + "\n"
		}
		return fmt.Errorf("errors occurred while deleting files:\n%s", errMsg)
	}

	g.Log().Infof(ctx, "成功删除图片:%d张", len(paths))
	return nil
}

// GetImgSrcFromStr 从给定字符串中提取图片路径信息
func GetImgSrcFromStr(input string) []string {
	// 定义匹配图片路径的正则表达式
	re := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

	// 查找所有匹配的图片路径
	matches := re.FindAllStringSubmatch(input, -1)

	// 预分配切片的容量
	imageSrcList := make([]string, 0, len(matches))

	for _, match := range matches {
		imageSrcList = append(imageSrcList, match[1])
	}

	return imageSrcList
}

// CheckFileAndPermissions 检查文件及权限问题
func CheckFileAndPermissions(path string) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 文件不存在直接忽略
		log.Printf("file does not exist: %+v", err)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error accessing file: %s, %v", path, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", path)
	}
	return nil
}
