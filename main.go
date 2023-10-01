package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

/*
	____   _                  _             __     __
	|  _ \ | |__   _   _  ___ (_)  ___  ___  \ \   / /___   _   _   __ _   __ _   ___  _ __
	| |_) || '_ \ | | | |/ __|| | / __|/ __|  \ \ / // _ \ | | | | / _` | / _` | / _ \| '__|
	|  __/ | | | || |_| |\__ \| || (__ \__ \   \ V /| (_) || |_| || (_| || (_| ||  __/| |
	|_|    |_| |_| \__, ||___/|_| \___||___/    \_/  \___/  \__, | \__,_| \__, | \___||_|
                   |___/                                    |___/         |___/
	Author: Physicser
	Date: 9.30.2023
*/

const version = "Voyager1.0"

// 输出 Logo
func printLogo() {
	var logo = "  ____   _                  _             __     __                                      \n" +
		" |  _ \\ | |__   _   _  ___ (_)  ___  ___  \\ \\   / /___   _   _   __ _   __ _   ___  _ __ \n" +
		" | |_) || '_ \\ | | | |/ __|| | / __|/ __|  \\ \\ / // _ \\ | | | | / _` | / _` | / _ \\| '__|\n" +
		" |  __/ | | | || |_| |\\__ \\| || (__ \\__ \\   \\ V /| (_) || |_| || (_| || (_| ||  __/| |   \n" +
		" |_|    |_| |_| \\__, ||___/|_| \\___||___/    \\_/  \\___/  \\__, | \\__,_| \\__, | \\___||_|   \n" +
		"                |___/                                    |___/         |___/         "

	fmt.Println(logo)

	fmt.Println("\n当前版本：", version, "\n")

	return
}

// 选择获取模式
func choseMode() {
	// 选择
	fmt.Println("请选择您要获取的内容\n" +
		"1. 账密（Account）2. 曲奇（Cookie）：")
	var getMode int
	_, err := fmt.Scanln(&getMode)

	// 错误处理
	if err != nil {
		fmt.Println("请您输入一个有效的整数！\n")
		time.Sleep(1 * time.Second)
		choseMode()
	}

	// 执行
	if getMode == 1 {
		getAccount()
	} else if getMode == 2 {
		getCookie()
	} else {
		fmt.Println("您选择的模式不存在，请重新选择！\n")
		time.Sleep(1 * time.Second)
		choseMode()
	}

	return
}

// Jsmh 人机验证
func verify(mode int) {
	// 请求响应内容
	response, err := http.Get("https://4399.js.mcdds.cn/captcha.php")
	if err != nil {
		fmt.Println("请求失败：", err, "\n")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭失败：", err, "\n")
			return
		}
	}(response.Body)

	// 读取响应内容
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("无法读取响应内容：%s\n", err)
		return
	}

	// 将内容保存到本地文件
	filename := "captcha.png"
	err = os.WriteFile(filename, body, 0644)
	if err != nil {
		fmt.Printf("保存验证码至本地文件失败：%s\n", err)
		return
	}

	fmt.Printf("\n已成功保存验证码至本地文件：%s\n", filename)

	// 打开验证码图片
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		// Linux （仅支持X11桌面环境）
		cmd = exec.Command("xdg-open", filename)
	case "darwin":
		// macOS
		cmd = exec.Command("open", filename)
	case "windows":
		// Windows
		cmd = exec.Command("cmd", "/C", "start", filename)
	default:
		log.Println("未知操作系统。请手动打开位于程序根目录的文件！")
	}
	if err := cmd.Run(); err != nil {
		log.Println(err, "文件打开失败。请手动打开位于程序根目录的文件！")
	}

	// 验证
	var captcha string
	fmt.Println("请输入验证码内容：")
	_, err = fmt.Scanln(&captcha)
	if err != nil {
		fmt.Println("未知错误\n")
		return
	}
	response, err = http.Get("https://4399.js.mcdds.cn/captcha_check.php?code=" + captcha)
	if err != nil {
		fmt.Println("请求失败：", err, "\n")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)
	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err, "\n")
		return
	}

	// 解析返回内容
	var verifyResults map[string]interface{}
	err = json.Unmarshal(body, &verifyResults)
	if err != nil {
		fmt.Println("Json解析失败：", err, "\n")
		return
	}

	// 后续操作
	message := verifyResults["msg"].(string)
	if message == "验证码正确" {
		fmt.Println("验证成功\n")
	}
	if mode == 1 {
		getAccount()
	} else if mode == 2 {
		getCookie()
	}

	return
}

// 获取账号
func getAccount() {
	// 请求返回内容
	response, err := http.Get("https://4399.js.mcdds.cn/get.php")
	if err != nil {
		fmt.Println("请求失败：", err, "\n")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err, "\n")
		return
	}

	// 解析返回内容
	var accountPasswd map[string]interface{}
	err = json.Unmarshal(body, &accountPasswd)
	if err != nil {
		fmt.Println("Json解析失败：", err, "\n")
		return
	}

	message := accountPasswd["msg"].(string)

	if message == "验证过期" || message == "请先验证" {
		verify(1)
		return
	}

	// 判断是否处于冷却时间
	cooldown := accountPasswd["code"].(float64)
	if cooldown == 1 {
		accountPasswd = accountPasswd["account"].(map[string]interface{})
		account := accountPasswd["account"].(string)
		passwd := accountPasswd["password"].(string)

		fmt.Println("账号：", account, "\n密码：", passwd, "\n")
	} else if cooldown == -4 {
		fmt.Println("服务冷却中\n")
	} else {
		fmt.Println("未知错误。错误代码：", cooldown, "\n")
	}

	return
}

func getCookie() {
	// 请求返回内容
	response, err := http.Get("https://4399.js.mcdds.cn/get_sauth.php")
	if err != nil {
		fmt.Println("请求失败：", err, "\n")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭失败：", err, "\n")
			return
		}
	}(response.Body)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败：", err, "\n")
		return
	}

	// 解析返回内容
	var cookies map[string]interface{}
	err = json.Unmarshal(body, &cookies)
	if err != nil {
		fmt.Println("Json解析失败：", err, "\n")
		return
	}

	// （验证）
	message := cookies["msg"].(string)
	if message == "验证过期" || message == "请先验证" {
		verify(2)
		return
	}

	// 输出
	cooldown := cookies["code"].(float64)
	if cooldown == 1 {
		cookie := cookies["account"].(string)
		fmt.Println(cookie, "\n")
	} else if cooldown == -4 {
		fmt.Println("服务冷却中\n")
	} else {
		fmt.Println("未知错误。错误代码：", cooldown, "\n")
	}

	return
}

func main() {
	printLogo()
	for {
		choseMode()
	}
}
