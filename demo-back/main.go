package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// 环境变量
var (
	staticPath string
	staticUrlPath string
	server = "0.0.0.0:5000"
	greeting = "你好"
)

// 返回码
const (
	RES_ERROR_UNKNOW = iota - 1
	RES_SUCCESS
	RES_ERROR_PARAMETER
	RES_ERROR_UPLOAD

)

func init() {
	// 文件服务器的服务路径
	if staticPathStr := os.Getenv("ENV_STATIC_PATH"); len(staticPath) > 0 {
		staticPath = staticPathStr
	} else {
		defaultStaticPath, err := filepath.Abs("media")
		if err != nil {
			panic(fmt.Sprintf("ERROR: static path init error: %s", err))
		}
		staticPath = defaultStaticPath
	}
	// 文件服务器的访问路径
	if staticUrlPathStr := os.Getenv("ENV_STATIC_URL_PATH"); len(staticUrlPathStr) > 0 {
		staticUrlPath = staticUrlPathStr
	} else {
		staticUrlPath = "media"
	}
	// 服务器访问地址
	if serverStr := os.Getenv("ENV_SERVER"); len(serverStr) > 0 {
		server = serverStr
	}
	// 问候语
	if greetingStr := os.Getenv("ENV_GREETING"); len(greetingStr) > 0 {
		greeting = greetingStr
	}
	// 初始化静态文件夹
	if _, err := os.Stat(staticPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(staticPath, os.ModePerm); err != nil {
				panic(fmt.Sprintf("ERROR: mkdir static path failed, error: %s", err))
			}
		}
	}
}

func main() {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.Use(Recover())



	v1Group := engine.Group("/v1")
	{
		v1Group.GET("/greeting", Greeting)
		v1Group.POST("/upload", Upload)
	}

	engine.Static(staticUrlPath, staticPath)
	engine.Run(server)
}

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ResponseData(c, RES_ERROR_UNKNOW, nil, err.(error))
			}
		}()
		c.Next()
	}
}

func ResponseData(c *gin.Context, code int, data interface{}, err error) {
	message := "Success"
	if err != nil {
		log.Println("ERROR: ", err)
		message = err.Error()
	}
	if code < 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
			Code: code,
			Message: err.Error(),
			Data: nil,
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, Response{
		Code: code,
		Message: message,
		Data: data,
	})
}

func Greeting(c *gin.Context) {
	i := 0
	fmt.Println(9/i)
	ResponseData(c, RES_SUCCESS, greeting, nil)
}

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		ResponseData(c, RES_ERROR_PARAMETER, nil, err)
		return
	}
	if err := c.SaveUploadedFile(file, path.Join(staticPath, file.Filename)); err != nil {
		ResponseData(c, RES_ERROR_UPLOAD, nil, err)
		return
	}
	ResponseData(c, RES_SUCCESS, path.Join("/"+staticUrlPath, file.Filename), nil)
}
