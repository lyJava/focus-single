package captcha

import (
	"context"
	"focus-single/internal/service"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/mojocn/base64Captcha"
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
	request := g.RequestFromCtx(ctx)
	captcha := base64Captcha.NewCaptcha(captchaDriver, captchaStore)

	_, content, answer := captcha.Driver.GenerateIdQuestionAnswer()
	item, _ := captcha.Driver.DrawCaptcha(content)
	g.Log().Infof(ctx, "验证码内容===%s", content)
	captchaStoreKey := guid.S()
	if err := request.Session.Set(name, captchaStoreKey); err != nil {
		g.Log().Errorf(ctx, "设置验证码session错误==%+v", err)
		return err
	}
	if err := captcha.Store.Set(captchaStoreKey, answer); err != nil {
		g.Log().Errorf(ctx, "设置验证码Store错误==%+v", err)
		return err
	}
	_, err := item.WriteTo(request.Response.Writer)
	return err
}

// VerifyAndClear 校验验证码，并清空缓存的验证码信息
func (s *sCaptcha) VerifyAndClear(r *ghttp.Request, name string, value string) bool {
	ctx := context.Background()
	captchaStoreKey := r.Session.MustGet(name).String()

	g.Log().Infof(ctx, "验证码:%s,前端传入验证码:%s,验证码缓存key==%s", name, value, captchaStoreKey)

	defer func(Session *gsession.Session, keys ...string) {
		err := Session.Remove(keys...)
		if err != nil {
			g.Log().Errorf(ctx, "验证码清空缓存错误==%+v", err)
		}
	}(r.Session, name)
	return captchaStore.Verify(captchaStoreKey, value, true)
}
