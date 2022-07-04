package utils

// your database msg
const ADDR, PORT, USER, PASSWORD, DATABASE = "ip", "3306", "user", "password", "your-database"

//const CLOUD_ELASTIC_SEARCH_HOST = "http://82.156.217.192:9200"

const HUAWEI_ELASTIC_SEARCH_HOST = "http://124.70.95.61:9200"

// 个人内网穿透ES使用
//const ELASTIC_SEARCH_HOST = "http://4r244274i3.zicp.vip:23656"

const LOCAL_ELASTIC_SEARCH_HOST = "http://172.18.0.1:9200"

// your elasticsearch host depand on your os type
var ELASTIC_SEARCH_HOST = HUAWEI_ELASTIC_SEARCH_HOST

const FRONTEND_URL = "https://slime.matrix53.top"

// your redis host
const REDIS_HOST = "124.70.95.61:6379"

// your redis password
const REDIS_PASSWORD = "redis1921@"

const EMAIL_HOST = "smtp.126.com"
const EMAIL_USER = ""
const EMAIL_PASSWORD = ""
const EMAIL_PORT = "465"
