package protocols

const (
	CodeOk = 0
	CodeError = 1
	CodeAccountExists    = 20 // 账号已存在
	CodeAccountNotExists = 21 // 账号不存在 
	CodeNotLogin         = 22
	CodeInvalidParameter         = 23 
)

func CodeMessage(code int) string {
	switch code {
	case CodeAccountExists:
		return "账号已存在"
	case CodeAccountNotExists:
		return "账号不存在"
	case CodeNotLogin:
		return "尚未登录" 
	}
	return "" // 未知错误
}
