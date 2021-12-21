package initialize

import (
	"fmt"
	"gitee.com/online-publish/slime-scholar-go/service"
	"log"
	"os"
	"time"

	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/utils"
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
		utils.ELASTIC_SEARCH_HOST = "http://172.18.0.1:9200"
		utils.BACK_PATH = "https://slime.matrix53.top/api/v1/upload"
		utils.UPLOAD_PATH = "/share/"

	} else if utils.SysType == "windows" {
		utils.LOG_FILE_PATH = "./"
		utils.ELASTIC_SEARCH_HOST = "http://124.70.95.61:9200"

	}
}

func Init() {
	InitElasticSearch()
	InitRedis()
	InitMySQL()
	InitOS()
}
