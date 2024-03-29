package initialize

import (
	"fmt"
	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"log"
	"os"
	"time"

	"github.com/BFlameSwift/SlimeScholar-Go/global"
	"github.com/BFlameSwift/SlimeScholar-Go/model"
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitMySQL 初始化 MySQL 的相关配置
func InitMySQL() {
	addr, port, username, password, database := utils.ADDR, utils.PORT, utils.USER, utils.PASSWORD, utils.DATABASE
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, addr, port, database)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // Disable color
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	global.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.User{},
		&model.Author{},
		&model.SubmitScholar{},
		&model.Affiliation{},
		&model.AuthorConnection{},
		&model.Comment{},
		&model.Conference{},
		&model.Follow{},
		&model.Like{},
		&model.Message{},
		&model.Paper{},
		&model.BrowsingHistory{},
		&model.Tag{},
		&model.TagPaper{},
		&model.Transfer{},

		// 生成新的数据库表
		&model.PaperReference{},
		&model.Collect{},
	)
}

// InitElasticSearch 初始化Elasticsearch 链接
func InitElasticSearch() {
	service.Init()
}

// InitRedis 初始化Redis连接
func InitRedis() {
	service.InitRedis()
}

// InitOS 根据OS的不同配置不同的变量
func InitOS() {
	if utils.SysType == "linux" {
		utils.LOG_FILE_PATH = "/backend/"
		// local host 内网环境下访问es 速度更快
		utils.ELASTIC_SEARCH_HOST = "http://172.18.0.1:9200"
		utils.BACK_PATH = "https://<your domain name>/api/v1/upload"
		// 服务器docker中挂载至/share/ 中
		utils.UPLOAD_PATH = "/share/"

	} else if utils.SysType == "windows" {
		utils.LOG_FILE_PATH = "./"
		// 本地环境访问es
		utils.ELASTIC_SEARCH_HOST = "http://127.0.0.1:9200"

	}
}

func Init() {
	InitElasticSearch()
	InitRedis()
	InitMySQL()
	InitOS()
}
