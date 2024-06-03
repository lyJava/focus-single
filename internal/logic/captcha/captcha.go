package captcha

import (
	"context"
	"focus-single/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/mojocn/base64Captcha"
	"log"
)

type sCaptcha struct{}

var (
	captchaStore  = base64Captcha.DefaultMemStore
	captchaDriver = newDriver()
)

func init() {
	service.RegisterCaptcha(New())
}

// New 验证码管理服务
func New() *sCaptcha {
	return &sCaptcha{}
}

func newDriver() *base64Captcha.DriverString {
	driver := &base64Captcha.DriverString{
		Height:          44,
		Width:           130,
		NoiseCount:      2,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowHollowLine,
		Length:          4,
		Source:          "123456789",
		Fonts:           []string{"wqy-microhei.ttc"},
	}
	return driver.ConvertFonts()
}

// NewAndStore 创建验证码，直接输出验证码图片内容到HTTP Response.
func (s *sCaptcha) NewAndStore(ctx context.Context, name string) error {
	captcha := base64Captcha.NewCaptcha(captchaDriver, captchaStore)
	id, content, answer := captcha.Driver.GenerateIdQuestionAnswer()
	log.Printf("验证码ID:%s,验证码:%s,结果:%s", id, content, answer)

	captchaStoreKey := guid.S()
	log.Printf("验证码缓存key:%s,值为:%s", name, captchaStoreKey)

	request := g.RequestFromCtx(ctx)
	if err := request.Session.Set(name, captchaStoreKey); err != nil {
		g.Log().Errorf(ctx, "设置验证码session错误==%+v", err)
		return err
	}

	if err := captcha.Store.Set(captchaStoreKey, answer); err != nil {
		g.Log().Errorf(ctx, "设置验证码Store错误==%+v", err)
		return err
	}

	item, _ := captcha.Driver.DrawCaptcha(content)
	_, err := item.WriteTo(request.Response.Writer)

	data, _ := request.Session.Data()
	log.Printf("创建验证码session中data===%v", data) // 这里r.Session.Data()有值的
	return err
}

// VerifyAndClear 校验验证码，并清空缓存的验证码信息
func (s *sCaptcha) VerifyAndClear(request *ghttp.Request, name string, value string) bool {
	defer func(Session *gsession.Session, keys ...string) {
		err := Session.Remove(keys...)
		if err != nil {

		}
	}(request.Session, name)

	dataMap, _ := request.Session.Data()
	log.Printf("获取验证码session中data===%v", dataMap) //这里的r.Session.Data()有时候为空

	captchaStoreKey := request.Session.MustGet(name).String()

	log.Printf("验证码验证:%s,传入验证码:%s,码缓存key==%s", name, value, captchaStoreKey)

	return captchaStore.Verify(captchaStoreKey, value, true)
}
