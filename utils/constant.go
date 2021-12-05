package utils

// token var
var (
	Secret     = "SlimeScholar" // 加盐
	ExpireTime = 3600 * 72      // token有效期
)

// error reason
const (
	ErrorServerBusy = "服务器繁忙"
	ErrorReLogin    = "请重新登陆"
)
