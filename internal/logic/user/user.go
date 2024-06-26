package user

import (
	"context"
	"fmt"
	"focus-single/internal/consts"
	"focus-single/internal/dao"
	"focus-single/internal/model"
	"focus-single/internal/model/do"
	"focus-single/internal/model/entity"
	"focus-single/internal/service"
	"focus-single/internal/util"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/o1egl/govatar"
	"log"
)

type sUser struct {
	avatarUploadPath      string // 头像上传路径
	avatarUploadUrlPrefix string // 头像上传对应的URL前缀
}

func init() {
	user := New()
	// 启动时创建头像存储目录
	if !gfile.Exists(user.avatarUploadPath) {
		if err := gfile.Mkdir(user.avatarUploadPath); err != nil {
			g.Log().Fatal(gctx.New(), err)
		}
	}
	service.RegisterUser(user)
}

func New() *sUser {
	return &sUser{
		avatarUploadPath:      g.Cfg().MustGet(gctx.New(), `upload.path`).String() + `/avatar`,
		avatarUploadUrlPrefix: `/upload/avatar`,
	}
}

// GetAvatarUploadPath 获得头像上传路径
func (s *sUser) GetAvatarUploadPath() string {
	return s.avatarUploadPath
}

// GetAvatarUploadUrlPrefix 获得头像上传对应的URL前缀
func (s *sUser) GetAvatarUploadUrlPrefix() string {
	return s.avatarUploadUrlPrefix
}

// Login 执行登录
func (s *sUser) Login(ctx context.Context, in model.UserLoginInput) error {
	userEntity, err := s.GetUserByPassportAndPassword(
		ctx,
		in.Passport,
		s.EncryptPassword(in.Passport, in.Password),
	)
	if err != nil {
		g.Log().Errorf(ctx, "用户登录错误===%+v", err)
		return err
	}
	if userEntity == nil {
		return gerror.New(`账号或密码错误`)
	}
	if err = service.Session().SetUser(ctx, userEntity); err != nil {
		g.Log().Errorf(ctx, "设置用户session错误===%+v", err)
		return err
	}
	// 自动更新上线
	service.BizCtx().SetUser(ctx, &model.ContextUser{
		Id:       userEntity.Id,
		Passport: userEntity.Passport,
		Nickname: userEntity.Nickname,
		Avatar:   userEntity.Avatar,
	})
	return nil
}

// Logout 注销
func (s *sUser) Logout(ctx context.Context) error {
	return service.Session().RemoveUser(ctx)
}

// EncryptPassword 将密码按照内部算法进行加密
func (s *sUser) EncryptPassword(passport, password string) string {
	return gmd5.MustEncrypt(passport + password)
}

// GetUserByPassportAndPassword 根据账号和密码查询用户信息，一般用于账号密码登录。
// 注意password参数传入的是按照相同加密算法加密过后的密码字符串。
func (s *sUser) GetUserByPassportAndPassword(ctx context.Context, passport, password string) (user *entity.User, err error) {
	err = dao.User.Ctx(ctx).Where(g.Map{
		dao.User.Columns().Passport: passport,
		dao.User.Columns().Password: password,
	}).Scan(&user)
	return
}

// CheckPassportUnique 检测给定的账号是否唯一
func (s *sUser) CheckPassportUnique(ctx context.Context, passport string) error {
	n, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Passport, passport).Count()
	if err != nil {
		g.Log().Errorf(ctx, "检测给定的账号错误===%+v", err)
		return err
	}
	if n > 0 {
		return gerror.Newf(`账号"%s"已被占用`, passport)
	}
	return nil
}

// CheckNicknameUnique 检测给定的昵称是否唯一
func (s *sUser) CheckNicknameUnique(ctx context.Context, nickname string) error {
	n, err := dao.User.Ctx(ctx).Where(dao.User.Columns().Nickname, nickname).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return gerror.Newf(`昵称"%s"已被占用`, nickname)
	}
	return nil
}

// Register 用户注册。
func (s *sUser) Register(ctx context.Context, in model.UserRegisterInput) error {
	return dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var user *entity.User
		if err := gconv.Struct(in, &user); err != nil {
			return err
		}
		if err := s.CheckPassportUnique(ctx, user.Passport); err != nil {
			return err
		}
		if err := s.CheckNicknameUnique(ctx, user.Nickname); err != nil {
			return err
		}
		user.Password = s.EncryptPassword(user.Passport, user.Password)
		// 自动生成头像
		avatarFilePath := fmt.Sprintf(`%s/%s.jpg`, s.avatarUploadPath, user.Passport)
		if err := govatar.GenerateFileForUsername(govatar.MALE, user.Passport, avatarFilePath); err != nil {
			g.Log().Errorf(ctx, "自动创建头像错误===%+v", err)
			return gerror.Wrapf(err, `自动创建头像失败`)
		}
		user.Avatar = fmt.Sprintf(`%s/%s.jpg`, s.avatarUploadUrlPrefix, user.Passport)
		log.Printf("注册用户:%s", util.ToJsonFormat(&user, false))
		_, err := dao.User.Ctx(ctx).TX(tx).
			Data(user).
			OmitEmpty().
			OnConflict(
				dao.User.Columns().Id,
			).
			Save()
		if err != nil {
			g.Log().Errorf(ctx, "用户注册错误===%+v", err)
			return err
		}
		return nil
	})
}

// UpdatePassword 修改个人密码
func (s *sUser) UpdatePassword(ctx context.Context, in model.UserPasswordInput) error {
	return dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		oldPassword := s.EncryptPassword(service.BizCtx().Get(ctx).User.Passport, in.OldPassword)
		n, err := dao.User.Ctx(ctx).TX(tx).
			Where(dao.User.Columns().Password, oldPassword).
			Where(dao.User.Columns().Id, service.BizCtx().Get(ctx).User.Id).
			Count()
		if err != nil {
			g.Log().Errorf(ctx, "修改密码错误===%+v", err)
			return err
		}
		if n == 0 {
			g.Log().Error(ctx, "原始密码错误")
			return gerror.New(`原始密码错误`)
		}
		newPassword := s.EncryptPassword(service.BizCtx().Get(ctx).User.Passport, in.NewPassword)
		_, err = dao.User.Ctx(ctx).Data(g.Map{
			dao.User.Columns().Password: newPassword,
		}).Where(dao.User.Columns().Id, service.BizCtx().Get(ctx).User.Id).Update()
		return err
	})
}

// GetProfileById 获取个人信息
func (s *sUser) GetProfileById(ctx context.Context, userId uint) (out *model.UserGetProfileOutput, err error) {
	if err = dao.User.Ctx(ctx).Where(dao.Content.Columns().Id, userId).Scan(&out); err != nil {
		return nil, err
	}
	// 需要判断nil是否存在,不存在需要判断为空,以防后续nil
	if out == nil {
		g.Log().Errorf(ctx, "获取个人信息错误===%+v", err)
		return nil, nil
	}
	out.Stats, err = s.GetUserStats(ctx, userId)
	if err != nil {
		g.Log().Errorf(ctx, "获取个人文章数量错误===%+v", err)
		return nil, err
	}
	return
}

// GetProfile 修改个人资料
func (s *sUser) GetProfile(ctx context.Context) (*model.UserGetProfileOutput, error) {
	return s.GetProfileById(ctx, service.BizCtx().Get(ctx).User.Id)
}

func (s *sUser) UpdateAvatar(ctx context.Context, in model.UserUpdateAvatarInput) error {
	return dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var (
			err error
		)
		_, err = dao.User.Ctx(ctx).TX(tx).
			Data(do.User{Avatar: in.Avatar}).
			Where(do.User{Id: in.UserId}).
			Update()
		if err != nil {
			g.Log().Errorf(ctx, "修改头像错误===%+v", err)
		}
		return nil
	})
}

// UpdateProfile 修改个人资料
func (s *sUser) UpdateProfile(ctx context.Context, in model.UserUpdateProfileInput) error {
	return dao.User.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var (
			err    error
			user   = service.BizCtx().Get(ctx).User
			userId = user.Id
		)
		n, err := dao.User.Ctx(ctx).TX(tx).
			Where(dao.User.Columns().Nickname, in.Nickname).
			WhereNot(dao.User.Columns().Id, userId).
			Count()
		if err != nil {
			g.Log().Errorf(ctx, "查询昵称数量错误===%+v", err)
			return err
		}
		if n > 0 {
			g.Log().Errorf(ctx, "昵称：%s已被占用", in.Nickname)
			return gerror.Newf(`昵称"%s"已被占用`, in.Nickname)
		}

		// .OmitEmpty()会过滤空值及默认值
		// _, err = dao.User.Ctx(ctx).TX(tx).OmitEmpty().Data(in).Where(dao.User.Columns().Id, userId).Update()
		whereMap := g.Map{}
		if in.Nickname != "" {
			whereMap[dao.User.Columns().Nickname] = in.Nickname
		}
		if in.Gender != 0 {
			whereMap[dao.User.Columns().Gender] = in.Gender
		}
		if in.Avatar != "" {
			whereMap[dao.User.Columns().Avatar] = in.Avatar
		}

		if len(whereMap) == 0 {
			return gerror.New(`没有更新的值，无法进行操作`)
		}

		_, err = dao.User.Ctx(ctx).TX(tx).
			Data(whereMap).
			Where(dao.User.Columns().Id, userId).
			Update()
		// 更新登录session Nickname
		if err == nil && user.Nickname != in.Nickname {
			sessionUser := service.Session().GetUser(ctx)
			sessionUser.Nickname = in.Nickname
			err = service.Session().SetUser(ctx, sessionUser)
		}
		if err != nil {
			g.Log().Errorf(ctx, "修改个人资料错误===%+v", err)
			return err
		}
		return nil
	})

}

// Disable 禁用指定用户
func (s *sUser) Disable(ctx context.Context, id uint) error {
	_, err := dao.User.Ctx(ctx).
		Data(dao.User.Columns().Status, consts.UserStatusDisabled).
		Where(dao.User.Columns().Id, id).
		Update()
	return err
}

// GetList 查询用户内容列表及用户信息
func (s *sUser) GetList(ctx context.Context, in model.UserGetContentListInput) (out *model.UserGetListOutput, err error) {
	out = &model.UserGetListOutput{}
	// 内容列表
	out.Content, err = service.Content().GetList(ctx, in.ContentGetListInput)
	if err != nil {
		g.Log().Errorf(ctx, "查询内容列表错误===%+v", err)
		return out, err
	}
	// 用户信息
	out.User, err = service.User().GetProfileById(ctx, in.UserId)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户信息错误===%+v", err)
		return out, err
	}
	// 统计信息
	out.Stats, err = s.GetUserStats(ctx, in.UserId)
	if err != nil {
		g.Log().Errorf(ctx, "查询文章数量错误===%+v", err)
		return out, err
	}
	return
}

// GetMessageList 消息列表
func (s *sUser) GetMessageList(ctx context.Context, in model.UserGetMessageListInput) (out *model.UserGetMessageListOutput, err error) {
	out = &model.UserGetMessageListOutput{
		Page: in.Page,
		Size: in.Size,
	}
	var userId = service.BizCtx().Get(ctx).User.Id
	// 管理员看所有的
	if !s.IsAdmin(ctx, userId) {
		in.UserId = userId
	}

	replyList, err := service.Reply().GetList(ctx, model.ReplyGetListInput{
		Page:       in.Page,
		Size:       in.Size,
		TargetType: in.TargetType,
		TargetId:   in.TargetId,
		UserId:     in.UserId,
	})
	if err != nil {
		return nil, err
	}
	if replyList != nil {
		out.List = replyList.List
	}
	out.Stats, err = s.GetUserStats(ctx, userId)
	if err != nil {
		g.Log().Errorf(ctx, "查询用户文章数量错误===%+v", err)
		return nil, err
	}
	return
}

// GetUserStats 获取文章数量
func (s *sUser) GetUserStats(ctx context.Context, userId uint) (map[string]int, error) {
	// 文章统计
	m := dao.Content.Ctx(ctx).Fields(dao.Content.Columns().Type, "count(*) total")
	if userId > 0 && !s.IsAdmin(ctx, userId) {
		m = m.Where(dao.Content.Columns().UserId, userId)
	}
	statsModel := m.Group(dao.Content.Columns().Type)
	statsAll, err := statsModel.All()
	if err != nil {
		g.Log().Errorf(ctx, "查询文章错误===%+v", err)
		return nil, err
	}
	statsMap := make(map[string]int)
	for _, v := range statsAll {
		value := v["type"]
		v2 := v["total"]
		statsMap[value.String()] = v2.Int()
	}
	// 回复统计
	replyModel := dao.Reply.Ctx(ctx).Fields("count(*) total")
	if userId > 0 && !s.IsAdmin(ctx, userId) {
		replyModel = replyModel.Where(dao.Reply.Columns().UserId, userId)
	}
	record, err := replyModel.One()
	if err != nil {
		g.Log().Errorf(ctx, "查询文章回复错误===%+v", err)
		return nil, err
	}
	value := record["total"]
	statsMap["message"] = value.Int()

	return statsMap, nil
}

// IsCtxAdmin 当前用户是否管理员
func (s *sUser) IsCtxAdmin(ctx context.Context) bool {
	var ctxUser = service.BizCtx().Get(ctx).User
	return s.IsAdmin(ctx, ctxUser.Id)
}

// IsAdmin 判断给定用户是否管理员
func (s *sUser) IsAdmin(ctx context.Context, userId uint) bool {
	adminIds := g.Cfg().MustGet(ctx, "setting.adminIds").Uints()
	for _, adminId := range adminIds {
		if userId == adminId {
			return true
		}
	}
	return false
}
