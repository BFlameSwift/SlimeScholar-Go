package initialize

import (
	"fmt"
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

// 初始化 MySQL 的相关配置
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
		//&model.Author{},
		&model.Affiliation{},
		&model.AuthorConnection{},
		&model.Comment{},
		&model.Follow{},
		&model.Like{},
		&model.Message{},
		//&model.Paper{},
		&model.BrowsingHistory{},
		&model.Tag{},
		// 生成新的数据库表
	)
}
