package consts

const (
	UserStatusOk       = 0 // 用户状态正常
	UserStatusDisabled = 1 // 用户状态禁用
	UserGenderMale     = 1 // 性别: 男
	UserGenderFemale   = 2 // 性别: 女
	UserGenderUnknown  = 3 // 性别: 未知
	UserLoginUrl       = "/login"
)

func GetGenderByType(typeInt int) string {
	dataMap := map[int]string{
		UserGenderMale:    "男",
		UserGenderFemale:  "女",
		UserGenderUnknown: "未知",
	}
	if content, exists := dataMap[typeInt]; exists {
		return content
	}
	return "未设置"
}
