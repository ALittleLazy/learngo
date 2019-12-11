package main

import (
	"MyTool/menu"
	"bufio"
	"fmt"
	"os"
)

func mainMenu() {
	fmt.Println()
	fmt.Println("mssql:MSSQL管理")
}

func main() {
	fmt.Println("欢迎使用本工具，输入-help可以查看帮助。")
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
