# Slime Scholar backend



初始化包 ： `go mod tidy`,下载与生成依赖 
运行 go run ./main.go git






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


swagger : http://localhost:8080/swagger/index.html



