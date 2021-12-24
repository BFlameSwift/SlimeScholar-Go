package utils

import "runtime"

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

const FOLLOW_USER_PREFIX = "follow"
const BE_FOLLOWED_USER_PREFIX = "befollow"

// 操作系统类型linux/windows
const SysType = runtime.GOOS

var LOG_FILE_PATH = "./"

const LOG_FILE_NAME = "scholar.log"

var BACK_PATH = "http://82.156.217.192:8000/api/v1/upload"

var UPLOAD_PATH = "./media/"

const PICTURE = "https://img-1304418829.cos.ap-beijing.myqcloud.com/avatar-grey-bg.jpg"

// 匹配到10000000文档就返回
const TERMINATE_AFTER = 10000000
