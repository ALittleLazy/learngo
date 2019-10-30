package main

import (
	"MyTool/conndatabase/mssql"
	"bufio"
	"fmt"
	"os"
)

func main() {

	fmt.Println("欢迎使用本工具，输入-help可以查看帮助。")
	fmt.Println("目前支持mssql数据库连接检测。")

	scanner := bufio.NewScanner(os.Stdin)

loop:
	for scanner.Scan() {
		switch scanner.Text() {
		case "mssql":
			mssql.Menu()
		default:
			break loop
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
