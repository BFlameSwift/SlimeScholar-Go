package utils

import (
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendRegisterEmail(themail string, number int) {
	subject := "欢迎注册Slime学术成果分享平台"
	// 邮件正文
	mailTo := []string{
		themail,
	}
	body := "Hello,This is a email,这是你的注册码" + strconv.Itoa(number)
	err := SendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
		return
	}
	fmt.Println("sendRegisterEmail successfully")
	return
}
func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	// 具体信息可在secret.go 中填写
	//mailConn := map[string]string{
	//  "user": "xxx@163.com",
	//  "pass": "your password",
	//  "host": "smtp.163.com",
	//  "port": "465",
	//}

	mailConn := map[string]string{
		"user": EMAIL_USER,
		"pass": EMAIL_PASSWORD,
		"host": EMAIL_HOST,
		"port": EMAIL_PORT,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()

	m.SetHeader("From", m.FormatAddress(mailConn["user"], "Slime scholar")) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

func SendReplyEmail(themail string, paper_id string) {
	subject := "您在Slime学术成果分享平台的评论收到了一条回复"
	// 邮件正文
	mailTo := []string{
		themail,
	}
	//发送文献链接
	body := "您评论的文章链接为" + "https://slime.matrix53.top" + "/article?v=" + paper_id
	err := SendMail(mailTo, subject, body)
	if err != nil {
		log.Println(err)
		fmt.Println("send fail")
		return
	}
	fmt.Println("sendRegisterEmail successfully")
	return
}
