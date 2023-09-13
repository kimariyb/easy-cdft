package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*
EasyCDFT 是厦门大学 Kimariyb 开发的一款全自动批处理使用 Multiwfn 计算概念密度泛函理论 CDFT 各种量的 Go 语言程序

实现思路：
	1. 首先获取当前文件夹下的所有文件，找到符合用户在 config.ini 中设定的文件类型
	2. 批量用 Multiwfn 执行选中的文件
	3. 读取文件后，使用流实现目的

@Name: EasyCDFT
@Author: Kimariyb
@Institution: XiaMen University
@Data: 2023-09-12
*/

/*
展示开始界面，包括版权，作者，版本等信息
*/
func show() {
	// 获取当前文件的绝对路径
	filePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("获取文件路径失败:", err)
		return
	}

	// 获取最后修改时间的时间戳
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("获取文件信息失败:", err)
		return
	}

	modTime := fileInfo.ModTime()
	timestamp := modTime.Format("2006-Jan-02")

	fmt.Println("@Name: EasyCDFT")
	fmt.Println("@Version: v1.0.0, @Release date:", timestamp)
	fmt.Println("@Developer: Kimariyb, Ryan Hsiun")
	fmt.Println("@Address: XiaMen University, School of Electronic Science and Engineering")
	fmt.Println("@Website: https://github.com/kimariyb/easy-cdft")

	// 获取当前日期和时间
	now := time.Now().Format("Jan-02-2006, 15:04:05")

	// 输出版权信息和问候语
	fmt.Printf("(Copyright (C) 2023 Kimariyb. Currently timeline: %s)\n", now)
	fmt.Println()
}

/*
读取当前文件夹下的 config.ini 文件
*/
func readIniFile(filePath string) (*ini.File, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

/*
获取 config 文件的路径
*/
func getConfigFilePath() string {
	currentPath, _ := os.Getwd()
	configFileName := "config.ini"
	configFilePath := filepath.Join(currentPath, configFileName)
	return configFilePath
}

/*
得到符合 config 中 inputType 属性配置的文件类型的所有文件
*/
func getFilesByType(inputType string) ([]string, error) {
	currentPath, _ := os.Getwd()
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == "."+inputType {
			result = append(result, file.Name())
		}
	}

	return result, nil
}

func processFile(file string, multiwfnPath string, commandLines []string) {
	fmt.Println("Documents being processed:", file)

	// 创建一个命令对象，执行 Multiwfn 程序
	cmd := exec.Command(multiwfnPath, file)

	// 获取标准输入管道
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 启动命令
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting process:", err)
		return
	}

	// 向命令的标准输入写入指令
	for _, line := range commandLines {
		_, err := io.WriteString(stdin, line+"\n")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	// 关闭标准输入管道
	stdin.Close()

	// 等待命令执行完成
	err = cmd.Wait()

	// 在每次执行完 Multiwfn 后会在当前文件夹下生成一个 CDFT.txt 文件
	// 将该文件名修改为 ${name}-CDFT.txt 文件，${name} 为一个变量，代表 file 文件的名字
	// 生成 CDFT.txt 文件名
	outputFile := strings.TrimSuffix(file, filepath.Ext(file)) + "-CDFT.txt"
	// 重命名 CDFT.txt 文件
	err = os.Rename("CDFT.txt", outputFile)
	if err != nil {
		fmt.Println("Error renaming file:", err)
		return
	}

	fmt.Println("Renamed CDFT.txt to", outputFile)
}

func main() {
	// 展示 hello 页面
	show()
	//config.ini 文件路径，必须在运行程序的当前路径下，名字为 config.ini
	configFilePath := getConfigFilePath()
	// 读取 config 文件
	cfg, err := readIniFile(configFilePath)
	if err != nil {
		fmt.Printf("Failed to read INI file: %v\n", err)
		return
	}
	// 读取 ini 文件中的 section
	section, err := cfg.GetSection("")
	if err != nil {
		fmt.Printf("Failed to get section: %v\n", err)
		return
	}
	// 读取 inputType 和 multiwfnPath 两个配置
	inputType := section.Key("inputType").String()
	multiwfnPath := section.Key("multiwfnPath").String()
	// 同时打印 inputType 和 multiwfnPath
	fmt.Printf("inputType: %s\n", inputType)
	fmt.Printf("multiwfnPath: %s\n", multiwfnPath)
	// 根据 inputType 所配置的文件类型，得到当前文件下所有符合 inputType 类型的文件
	files, err := getFilesByType(inputType)
	if err != nil {
		fmt.Printf("Failed to get files: %v\n", err)
		return
	}
	fmt.Println("Files with inputType", inputType+":")
	for _, file := range files {
		fmt.Println(file)
	}

	//读取 mission 配置
	mission := section.Key("mission").String()
	// 读取 calcLevel 配置
	calcLevel := section.Key("calcLevel").String()
	// 读取 chargeSpin1、2、3 配置
	chargeSpin1 := section.Key("chargeSpin1").String()
	chargeSpin2 := section.Key("chargeSpin2").String()
	chargeSpin3 := section.Key("chargeSpin3").String()
	// 定义需要实现的流，创建一个切片用于存放命令
	var commandLines []string
	// 根据 config.ini 设置储存命令
	if mission == "0" {
		// 如果 mission 为 0 则使用计算指数命令
		commandLines = append(commandLines, "22", "1", calcLevel, chargeSpin1, chargeSpin2, chargeSpin3, "y", "2", "q")
	} else if mission == "1" {
		// 如果 mission 为 1 则使用计算福井函数命令
		fmt.Println("Warning: This feature is not yet developed")
		return
	} else if mission == "2" {
		// 如果 mission 为 2 则使用
		fmt.Println("Warning: This feature is not yet developed")
		return
	} else {
		// 如果没有符合的命令，则抛出异常
		fmt.Println("Error: Invalid mission value")
		// 返回错误或中止程序
		return
	}

	// 批量执行 Multiwfn
	for _, file := range files {
		processFile(file, multiwfnPath, commandLines)
	}

	fmt.Println()
	fmt.Println("The mission has been successfully completed!")

}
