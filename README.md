# Slime Scholar backend



初始化包 ： `go mod tidy`,下载与生成依赖 

运行 go run ./main.go git


## 项目结构

**docs** 
    不需要管，都由swag init 自动生成

**api/v1**
    主要放置Api结构，其中在v1 包下都可统一调用与不同文件无关，可以方便区分不同api的归属

**global**
    不用管，设置全局变量DB

**initialize**
    初始化项目使用，分为两个文件夹
* gorm.go 主要用于初始化与链接数据库
* 设置后端基本的全局路径 : /api/v1

**middleware**
    中间件，但是我还没搞太懂（，豪哥说别管
    设置一些允许的域名 与方法之类的吧

**model**
    放置所有的数据模型，相关语法在其中，会自动的生成数据

**router**
    每次写完了一个api就在这里配置一下router即可

**service**
    服务于api的函数，（不用过多解释吧

**test**
    用于测试一些函数的使用，正式项目可删除

**utils**
    工具类，放一些常量、配置之类的函数


**main.go**
    主入口， go run main.go 即可运行整个项目

**go mod**
    目前我的理解是go的包管理工具，使用前要用 go mod tidy 来生成与查找依赖包，很好用





### Swagger

[下载地址](https://github.com/go-swagger/go-swagger/releases/tag/v0.27.0), 最好加到环境变量中

[推荐入门链接](https://www.jianshu.com/p/4875b5ac9feb)

生成文档 `swag init `
    * 前置条件
        设置好GOROOT、GOPATH 把GOPATH/bin 也加到环境变量

        go get -u github.com/swaggo/swag/cmd/swag

        go get -u github.com/swaggo/gin-swagger

        go install  github.com/swaggo/swag/cmd/swag

        go install github.com/swaggo/gin-swagger



swagger : http://localhost:8000/swagger/index.html



