package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Path  string             `yaml:"path"`
	Main  RepositoryConfig   `yaml:"main"`
	Other []RepositoryConfig `yaml:"other"`
}
type RepositoryConfig struct {
	Repository string `yaml:"repository"`
	Branch     string `yaml:"branch"`
	Project    string `yaml:"project"`
}

type Repository struct {
	Name string
}

type GithubJson struct {
	Repository Repository
	Ref        string
}

var logger *zap.SugaredLogger

func GetConfig() (*Config, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	var data Config
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}
	config = &data
	return &data, nil
}

var config *Config

// main
//
//	@Description: 思路：收到github到webhook后，从环境变量获取仓库，分支。将代码clone到本地，然后go run 启动
func main() {
	//  初始化配置
	_, err := GetConfig()
	if err != nil {
		panic(err)
		return
	}
	// 初始化 Zap 日志记录器
	loggerService, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	logger = loggerService.Sugar()
	defer loggerService.Sync()
	//  初始化git项目
	err = updateGits()
	if err != nil {
		panic(err)
	}
	// webhook服务
	r := gin.Default()
	r.POST("/webhook", func(c *gin.Context) {
		data := GithubJson{}
		err := c.ShouldBindJSON(&data)
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		if data.Ref != config.Main.Branch && data.Ref != fmt.Sprintf("refs/heads/%v", config.Main.Branch) {
			err := errors.New(fmt.Sprintf("data.Ref=%v,branch=%v", data.Ref, config.Main.Branch))
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		err = updateGits()
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		err = startServer(config.Main.Project, config.Path)
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	})
	r.GET("/test", func(c *gin.Context) {
		err = updateGits()
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		err = startServer(config.Main.Project, config.Path)
		if err != nil {
			logger.Error(err)
			c.JSON(400, gin.H{"data": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": "success"})
	})
	_ = r.Run(fmt.Sprintf(":8000"))
}

// updateGits
//
//	@Description: 更新所有仓库
//	@return error
func updateGits() error {
	err := updateGit(config.Main.Project, config.Main.Repository, config.Main.Branch, config.Path)
	if err != nil {
		return err
	}
	for _, repositoryConfig := range config.Other {
		err = updateGit(repositoryConfig.Project, repositoryConfig.Repository, repositoryConfig.Branch, config.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

// updateGit
//
//	@Description: 更新单个仓库
//	@param project
//	@param repository
//	@param branch
//	@param path
//	@return error
func updateGit(project, repository, branch, path string) error {
	// 获取shell文件
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	shell := fmt.Sprintf("%v/git.sh", wd)
	//@仓库地址 @分支 @本地路径 @项目文件夹
	shellCmd := fmt.Sprintf("%v '%v' '%v' '%v' '%v'", shell, repository, branch, path, project)
	return execShell(shellCmd)
}

// startServer
//
//	@Description: 启动服务
//	@param project
//	@param path
//	@return error
func startServer(project, path string) error {
	// 获取shell文件
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	shell := fmt.Sprintf("%v/start.sh", wd)
	// 参数：@项目名称 project @路径 path
	shellCmd := fmt.Sprintf("%v '%v' '%v'", shell, project, path)
	return execShell(shellCmd)
}

// stopServer
//
//	@Description: 停止服务
//	@param project
//	@return error
func stopServer(project string) error {
	// 获取shell文件
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	shell := fmt.Sprintf("%v/stop.sh", wd)
	// 参数：@项目名称 project
	shellCmd := fmt.Sprintf("%v '%v'", shell, project)
	return execShell(shellCmd)
}

// execShell
//
//	@Description: 执行shell命令
//	@param shell
//	@return error
func execShell(shell string) error {
	err := os.Chmod(shell, 0755)
	if err != nil {
		return err
	}
	// 创建一个具有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建命令对象
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", shell)
	output, err := cmd.CombinedOutput()
	logger.Infof("output=%v", string(output))
	if err != nil {
		// 如果命令因超时而终止，则捕获 context.DeadlineExceeded 错误
		if ctx.Err() == context.DeadlineExceeded {
			logger.Error("命令执行超时")
			return err
		}
		logger.Errorf("命令执行失败:%v", err)
		return err
	}
	return nil
}
