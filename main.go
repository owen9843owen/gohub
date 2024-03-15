package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"time"
)

type Repository struct {
	Name string
}

type GithubJson struct {
	Repository Repository
	Ref        string
}

var logger *zap.SugaredLogger

// main
//
//	@Description: 思路：收到github到webhook后，从环境变量获取仓库，分支。将代码clone到本地，然后go run 启动
func main() {
	// 初始化 Zap 日志记录器
	loggerService, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = loggerService.Sugar()
	defer loggerService.Sync()
	// 获取shell文件
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	shell := fmt.Sprintf("%v/goserver.sh", wd)
	// 环境变量
	serverPort := os.Getenv("port")                   // 端口
	serverName := os.Getenv("server")                 // 服务名称
	path := os.Getenv("path")                         // 路径
	commonRepository := os.Getenv("commonRepository") // 仓库名称
	repository := os.Getenv("repository")             // 仓库名称
	branch := os.Getenv("branch")                     // 分支

	r := gin.Default()

	// webhook监听
	r.POST("/webhook", func(c *gin.Context) {
		data := GithubJson{}
		err := c.ShouldBindJSON(&data)
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		if data.Repository.Name == repository {
			err := errors.New("仓库地址错误")
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		if data.Ref != branch {
			err := errors.New(fmt.Sprintf("data.Ref=%v,branch=%v", data.Ref, branch))
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		err = os.Chmod(shell, 0755)
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		cmd := fmt.Sprintf("%v %v %v %v %v %v", shell, serverName, commonRepository, path, repository, branch)
		logger.Infof("cmd: %v", cmd)
		out, err := exec.Command("/bin/bash", cmd).Output()
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Infof("shell result: %v", out)
		c.JSON(200, gin.H{"data": "success"})
	})
	r.GET("/test", func(c *gin.Context) {
		err = os.Chmod(shell, 0755)
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		shellCmd := fmt.Sprintf("%v '%v' '%v' '%v' '%v' '%v'", shell, serverName, commonRepository, path, repository, branch)
		logger.Infof("shellCmd: %v", shellCmd)
		err := execShell(shellCmd)
		//// 创建一个具有超时的上下文
		//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancel()
		//cmd := exec.CommandContext(ctx, "/bin/bash", "-c", shellCmd)
		//out, err := cmd.Output()
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	})
	_ = r.Run(fmt.Sprintf(":%v", serverPort))
}

func execShell(shell string) error {
	// 创建一个具有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建命令对象
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", shell)

	// 启动命令
	if err := cmd.Start(); err != nil {
		logger.Errorf("命令启动失败:%v", err)
		return err
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		// 如果命令因超时而终止，则捕获 context.DeadlineExceeded 错误
		if ctx.Err() == context.DeadlineExceeded {
			logger.Error("命令执行超时")
			return err
		}
		logger.Errorf("命令执行失败:%v", err)
		return err
	}
	output, err := cmd.Output()
	if err != nil {
		logger.Error(err)
		return err
	}
	logger.Infof("%v", string(output))
	return nil
}
