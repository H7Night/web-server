package errmsg

const (
	Success = 200
	Error   = 500

	ErrorUsernameUsed   = 1001
	ErrorPasswordWrong  = 1002
	ErrorUserNotExist   = 1003
	ErrorTokenExist     = 1004
	ErrorTokenRuntime   = 1005
	ErrorTokenWrong     = 1006
	ErrorTokenTypeWrong = 1007
	ErrorUserNoRight    = 1008
)

var codeMsg = map[int]string{
	Success:             "OK",
	Error:               "FAIL",
	ErrorUsernameUsed:   "用户名已存在！",
	ErrorPasswordWrong:  "密码错误",
	ErrorUserNotExist:   "用户不存在或已被删除",
	ErrorTokenExist:     "TOKEN不存在",
	ErrorTokenRuntime:   "TOKEN已过期",
	ErrorTokenWrong:     "TOKEN不正确",
	ErrorTokenTypeWrong: "TOKEN格式错误",
	ErrorUserNoRight:    "该用户无权限",
}

func GetErrMsg(code int) string {
	return codeMsg[code]
}
