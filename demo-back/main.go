package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	RES_SUCCESS = iota
	RES_ERROR_PARAMETER
	RES_ERROR_UPLOAD
	RES_ERROR_UNKNOW
)

func init() {
	// 文件服务器的服务路径
	if staticPathStr := os.Getenv("ENV_STATIC_PATH"); len(staticPath) > 0 {
		staticPath = staticPathStr
	} else {
		defaultStaticPath, err := filepath.Abs("static")
		if err != nil {
			panic(fmt.Sprintf("ERROR: static path init error: %s", err))
		}
		staticPath = defaultStaticPath
	}
	// 文件服务器的访问路径
	if staticUrlPathStr := os.Getenv("ENV_STATIC_URL_PATH"); len(staticUrlPathStr) > 0 {
		staticUrlPath = staticUrlPathStr
	} else {
		staticUrlPath = "static"
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
	engine := gin.Default()
	v1Group := engine.Group("/v1")
	{
		v1Group.GET("/greeting", Greeting)
	}
	engine.POST("/upload", Upload)
	engine.Static(staticUrlPath, staticPath)
	engine.Run(server)
}

type Response struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func Greeting(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: RES_SUCCESS,
		Message: "Success",
		Data: greeting,
	})
}

func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: RES_ERROR_PARAMETER,
			Message: err.Error(),
		})
		return
	}
	if err := c.SaveUploadedFile(file, path.Join(staticPath, file.Filename)); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code: RES_ERROR_UPLOAD,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code: RES_SUCCESS,
		Message: "Success",
		Data: path.Join("/"+staticUrlPath, file.Filename),
	})
}
