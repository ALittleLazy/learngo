package main

import (
	"MyTool/menu"
	"bufio"
	"fmt"
	"os"
)

func mainMenu() {
	fmt.Println()
	fmt.Println("MSSQL管理:mssql")
}

func main() {
	fmt.Println("欢迎使用本工具，目前支持数据连接及作业检查。")
	mainMenu()

	scanner := bufio.NewScanner(os.Stdin)

loop:
	for scanner.Scan() {
		switch scanner.Text() {
		case "mssql":
			menu.Menu()
		default:
			break loop
		}
		mainMenu()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
