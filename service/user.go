package service

// user 相关 service
import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"gitee.com/online-publish/slime-scholar-go/global"
	"gitee.com/online-publish/slime-scholar-go/model"
	"gitee.com/online-publish/slime-scholar-go/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

// 根据用户email 查询某个用户
func QueryAUserByEmail(email string) (user model.User, notFound bool) {
	err := global.DB.Where("email = ?", email).First(&user).Error
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
func UpdateConfirmAUser(user *model.User, has_comfirmed bool) error {
	if has_comfirmed == false {
		user.ConfirmNumber = rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
		utils.SendRegisterEmail(user.Email, user.ConfirmNumber)
		err := global.DB.Save(user).Error
		return err
	}
	user.HasConfirmed = true
	err := global.DB.Save(user).Error
	return err
}

// 发送猪猪邮件
func SendRegisterEmail(themail string, number int) {
	subject := "欢迎注册Slime学术成果分享平台"
	// 邮件正文
	mailTo := []string{
		themail,
	}
	body := "Hello,This is a email,这是你的注册码" + strconv.Itoa(number)
	err := utils.SendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
		return
	}
	fmt.Println("sendRegisterEmail successfully")
	return
}

// 获取Token
func GetToken(claims *model.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(utils.Secret))
	if err != nil {
		return "", errors.New(utils.ErrorServerBusy)
	}
	return signedToken, nil

}

// 验证token
func Verify(c *gin.Context) {
	strToken := c.Param("token")
	claim, err := VerifyAction(strToken)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, "verify,", claim.Username)
}

func VerifyAction(strToken string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(strToken, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Secret), nil
	})
	if err != nil {
		return nil, errors.New(utils.ErrorServerBusy)
	}
	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return nil, errors.New(utils.ErrorReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, errors.New(utils.ErrorServerBusy)
	}
	fmt.Println("verify")
	return claims, nil
}
func VerifyAuthorization(strToken string, userID uint64, username, password string) (bool, error) {
	token, err := jwt.ParseWithClaims(strToken, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(utils.Secret), nil
	})
	if err != nil {
		return false, errors.New(utils.ErrorServerBusy)
	}
	claims, ok := token.Claims.(*model.JWTClaims)
	if !ok {
		return false, errors.New(utils.ErrorReLogin)
	}
	if err := token.Claims.Valid(); err != nil {
		return false, errors.New(utils.ErrorServerBusy)
	}
	fmt.Println("verifying")
	if claims.UserID != userID || claims.Username != username {
		return false, nil
	}
	return true, nil

}
func CreateASubmit(submit *model.SubmitScholar) (err error) {
	if err = global.DB.Create(&submit).Error; err != nil {
		return err
	}
	return nil
}

// 根据提交Submit ID 查询某个Submit
func QueryASubmitByID(submit_id uint64) (submit model.SubmitScholar, notFound bool) {
	err := global.DB.Where("submit_id = ?", submit_id).First(&submit).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return submit, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return submit, false
	}
}
