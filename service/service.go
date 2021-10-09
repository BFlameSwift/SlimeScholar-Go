package service

import (
	"errors"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"math/rand"
	"time"

	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gorm.io/gorm"
)

// 创建用户
func CreateAUser(user *model.User) (err error) {
	if err = global.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// 根据用户 ID 查询某个用户
func QueryAUserByID(userID uint64) (user model.User, notFound bool) {
	err := global.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return user, false
	}
}

// 根据用户 username 查询某个用户
func QueryAUserByUsername(username string) (user model.User, notFound bool) {
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return user, false
	}
}

// 更新用户的用户名、密码、个人信息
func UpdateAUser(user *model.User, username string, password string, userInfo string) error {
	user.Username = username
	user.Password = password
	user.UserInfo = userInfo
	err := global.DB.Save(user).Error
	return err
}
// 如果bool == false 重发邮件，否则就把user的comfirm = true
func UpdateConfirmAUser(user *model.User, has_comfirmed bool) error{
	if has_comfirmed == false{
		user.ConfirmNumber = rand.New(rand.NewSource(time.Now().UnixNano())).Int()%1000000
		utils.SendRegisterEmail(user.Email,user.ConfirmNumber)
		err := global.DB.Save(user).Error
		return err
	}
	user.HasComfirmed = true
	err := global.DB.Save(user).Error
	return err
}
