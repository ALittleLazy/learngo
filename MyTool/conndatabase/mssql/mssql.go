package mssql

import (
	"MyTool/readwrite"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConnParamater struct{
	userid,password,ip,database,remark string
}

func ConnMssql(sc *bufio.Scanner){
	var connParaSet ConnParamater

	fmt.Print("ip address:")
	sc.Scan()
	connParaSet.ip = sc.Text()

	fmt.Print("database name:")
	sc.Scan()
	connParaSet.database = sc.Text()

	fmt.Print("conn remark:")
	sc.Scan()
	connParaSet.remark = sc.Text()

	fmt.Print("use id:")
	sc.Scan()
	connParaSet.userid = sc.Text()

	fmt.Print("password:")
	sc.Scan()
	connParaSet.password = sc.Text()

	connstr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable",connParaSet.ip,connParaSet.userid,connParaSet.password,connParaSet.database)

	db,err := sql.Open("mssql",connstr)
	if err != nil{
		log.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil{
		log.Println("数据连接失败：",err)
		return
	}

	fmt.Println(connParaSet.ip,"连接成功")

	str := fmt.Sprintf("%s|%s|%s|%s|%s|\n",connParaSet.ip,connParaSet.database,connParaSet.userid,connParaSet.password,connParaSet.remark)
	file,_ := exec.LookPath(os.Args[0])
	file_path,_ := filepath.Abs(file)
	file_path,filename := filepath.Split(file_path)
	filename = strings.Replace(filename,filepath.Ext(filename),".ini",-1)
	if err = readwrite.WriteString(str,filepath.Join(file_path,filename));err != nil{
		log.Println(err)
	}
}

