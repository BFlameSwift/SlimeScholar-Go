package v1

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/BFlameSwift/SlimeScholar-Go/model"
	"github.com/BFlameSwift/SlimeScholar-Go/service"
	"github.com/BFlameSwift/SlimeScholar-Go/utils"
	"github.com/gin-gonic/gin"
)

// Register doc
// @description 注册
// @Tags 用户管理
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string false "用户邮箱"
// @Success 200 {string} string "{"success": true, "message": "用户创建成功"}"
// @Failure 200 {string} string "{"success": false, "message": "用户已存在"}"
// @Router /user/register [POST]
func Register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	email := c.Request.FormValue("email")
	userType, _ := strconv.ParseUint(c.Request.FormValue("user_type"), 0, 64)
	userInfo, affiliation, userType := "", "", 0
	user_confirm_number := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	//affiliation := c.Request.FormValue("affiliation")
	user := model.User{Username: username, Password: password, UserInfo: userInfo, UserType: userType, Affiliation: affiliation, Email: email, ConfirmNumber: user_confirm_number, RegTime: time.Now()}
	_, notFound := service.QueryAUserByUsername(username)
	if notFound {
		service.CreateAUser(&user)
		tag := model.Tag{TagName: "默认", UserID: user.UserID, CreateTime: time.Now(), Username: user.Username}
		service.CreateATag(&tag)
		utils.SendRegisterEmail(email, user.ConfirmNumber)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户创建成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户已存在"})
	}
}

// Confirm doc
// @description 验证邮箱
// @Tags 用户管理
// @Param username formData string true "用户名"
// @Param confirm_number formData int true "confirm_number"
// @Success 200 {string} string "{"success": true, "message": "用户验证邮箱成功"}"
// @Failure 401 {string} string "{"success": false, "message": "用户已验证邮箱"}"
// @Failure 402 {string} string "{"success": false, "message": "用户输入验证码错误}"
// @Failure 404 {string} string "{"success": false, "message": "用户不存在}"
// @Failure 600 {string} string "{"success": false, "message": "用户待修改，传入false 更新验证码，否则为验证正确}"
// @Router /user/confirm [POST]
func Confirm(c *gin.Context) {

	confirm_number := c.Request.FormValue("confirm_number")

	username := c.Request.FormValue("username")
	user, notFound := service.QueryAUserByUsername(username)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户不存在", "status": 404})
	} else {
		if user.HasConfirmed == true {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户已验证", "status": 401})
		} else {
			num, _ := strconv.Atoi(confirm_number)
			if num != user.ConfirmNumber {
				c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户输入验证码错误", "status": 402})
			} else {
				service.UpdateConfirmAUser(&user, true)
				c.JSON(http.StatusOK, gin.H{"success": true, "message": "用户验证成功", "status": 200})
			}
		}
	}

}

// Login doc
// @description 登录
// @Tags 用户管理
// @Param username formData string false "用户名"
// @Param email formData string false "用户邮箱"
// @Param password formData string true "密码"
// @Success 200 {string} string "{"success": true, "message": "登录成功", "detail": user的信息}"
// @Failure 402 {string} string "{"success": false, "message": "密码错误"}"
// @Failure 401 {string} string "{"success": false, "message": "没有该用户"}"
// @Router /user/login [POST]
func Login(c *gin.Context) {
	username := c.Request.FormValue("username")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	user, notFound := model.User{}, true
	if username != "" {
		user, notFound = service.QueryAUserByUsername(username)
	} else {
		user, notFound = service.QueryAUserByEmail(email)
	}
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 401})
	} else {
		if user.Password != password {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": "密码错误", "status": 402})
		} else {
			if user.HasConfirmed == false {
				c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户尚未确认邮箱", "status": 403})
			} else {
				claims := &model.JWTClaims{
					UserID:   user.UserID,
					Username: user.Username,
					Password: password,
				}
				claims.IssuedAt = time.Now().Unix()
				claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(utils.ExpireTime)).Unix()
				signedToken, err := service.GetToken(claims)
				if err != nil {
					c.String(http.StatusNotFound, err.Error())
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"success":       true,
					"message":       "登录成功",
					"detail":        user,
					"status":        200,
					"Authorization": signedToken})
			}
		}
	}
}

// ModifyUser doc
// @description 修改用户信息（支持修改用户名和密码）
// @Tags 用户管理
// @Security Authorization
// @param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param username formData string true "用户名"
// @Param user_info formData string true "用户个人信息"
// @Param password_old formData string true "原密码"
// @Param password_new formData string true "新密码"
// @Success 200 {string} string "{"success": true, "message": "修改成功", "data": "model.User的所有信息"}"
// @Failure 400 {string} string "{"success": false, "message": "用户未登录"}"
// @Failure 401 {string} string "{"success": false, "message": "原密码输入错误"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Failure 500 {string} string "{"success": false, "message": "数据库操作时的其他错误"}"
// @Router /user/modify [POST]
func ModifyUser(c *gin.Context) {

	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	username := c.Request.FormValue("username")
	userInfo := c.Request.FormValue("user_info")
	passwordOld := c.Request.FormValue("password_old")
	passwordNew := c.Request.FormValue("password_new")
	user, notFoundUserByID := service.QueryAUserByID(userID)

	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return
	}

	authorization := c.Request.Header.Get("Authorization")
	verify_answer, _ := service.VerifyAuthorization(authorization, userID, username, user.Password)

	if authorization == "" || !verify_answer {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 400, "message": "用户未登录"})
		return
	}

	if passwordOld != user.Password {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  401,
			"message": "原密码输入错误",
		})
		return
	}
	_, notFoundUserByName := service.QueryAUserByUsername(username)
	if !notFoundUserByName && username != user.Username {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "用户名已被占用",
		})
		return
	}
	err := service.UpdateAUser(&user, username, passwordNew, userInfo)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	//data, _ := jsoniter.Marshal(&user)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功",
		"status":  200,
		"data":    user,
	})
}

// TellUserInfo doc
// @description 查看用户个人信息
// @Tags 用户管理
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "查看用户信息成功", "data": "model.User的所有信息"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Router /user/info [POST]
func TellUserInfo(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(userID)

	// @Param Authorization formData string false "Authorization"
	// authorization := c.Request.FormValue("Authorization")
	// authorization := c.Request.Header("Authorization")
	authorization := c.Request.Header.Get("Authorization")
	verify_answer, _ := service.VerifyAuthorization(authorization, userID, user.Username, user.Password)

	if authorization == "" || !verify_answer {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 400, "message": "用户未登录"})
		return
	}

	if notFoundUserByID {
		c.JSON(404, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看用户信息成功",
		"data":    user,
	})
}

type Img struct {
	Id     bson.ObjectId `bson:"_id"`
	ImgUrl string        `bson:"imgUrl"`
}

// ExportAvatar doc
// @description 上传头像
// @Tags 用户管理
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Param avatar formData file true "头像照片"
// @Success 200 {string} string "{"success": true, "message": "上传成功",}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Router /user/export/avatar [POST]
func ExportAvatar(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	//获取文件头
	imgFile, imgHead, imgErr := c.Request.FormFile("avatar")
	if imgErr != nil {
		fmt.Println(imgErr)
		return
	}
	defer imgFile.Close()

	imgFormat := strings.Split(imgHead.Filename, ".")
	var img Img
	img.Id = bson.NewObjectId()
	img.ImgUrl = user.Username + "_" + img.Id.Hex()[0:6] + "." + imgFormat[len(imgFormat)-1]

	image, e := os.Create(utils.UPLOAD_PATH + img.ImgUrl)
	if e != nil {
		fmt.Println(e)
		_ = os.Mkdir(utils.UPLOAD_PATH, 777)

		image, e = os.Create(utils.UPLOAD_PATH + img.ImgUrl)
		if e != nil {
			fmt.Println(e)
		}
		return
	}
	defer image.Close()

	_, e = io.Copy(image, imgFile)
	if e != nil {
		fmt.Println(e)
		return
	}

	errr := service.ExportAvatar(&user, "/upload/media/"+img.ImgUrl)
	if errr != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  500,
			"message": errr.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传成功",
		"status":  200,
		"data":    user.Avatar,
	})
}

// GetAvatar doc
// @description 获取头像
// @Tags 用户管理
// @Security Authorization
// @Param Authorization header string false "Authorization"
// @Param user_id formData string true "用户ID"
// @Success 200 {string} string "{"success": true, "message": "上传成功", "data" : "图像"}"
// @Failure 404 {string} string "{"success": false, "message": "用户ID不存在"}"
// @Router /user/get/avatar [POST]
func GetAvatar(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	authorization := c.Request.Header.Get("Authorization")
	user, err := VerifyLogin(userID, authorization, c)
	if err {
		return
	}

	var avatar string
	if user.Avatar == "" || len(user.Avatar) == 0 {
		avatar = utils.PICTURE
	} else {
		avatar = user.Avatar
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传成功",
		"status":  200,
		"data":    avatar,
	})
}
