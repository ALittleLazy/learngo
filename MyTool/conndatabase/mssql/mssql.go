package mssql

import (
	"MyTool/readwrite"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type connParamater struct {
	tp, ip, database, remark, userid, password string
}

//var databaseip map[string] connParamater

func Menu() {
	var file, file_path, filename string
	file, _ = exec.LookPath(os.Args[0])
	file_path, _ = filepath.Abs(file)
	file_path, filename = filepath.Split(file_path)
	filename = strings.Replace(filename, filepath.Ext(filename), ".ini", -1)
	//filename = "___go_build_CanDo_go.ini"

	menuno := make([]string, 10, 20)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("n:新建连接；r:显示现有链接；c:测试现有链接")

	databaseip, err := readINI(file, file_path, filename)
	if err != nil {
		log.Println(err)
	}

loop:
	for scanner.Scan() {
		switch scanner.Text() {
		case "n":
			newConnStr, err := newConnMssqlPing(scanner)
			if err != nil {
				log.Println(err)
			} else {
				databaseip[newConnStr.ip+"_"+newConnStr.database] = newConnStr
				if err := saveConnStr(databaseip, file, file_path, filename); err != nil {
					log.Println(err)
					continue
				}
			}
		case "r":
			i := 0
			for no, val := range databaseip {
				menuno[i] = no
				fmt.Println(i+1, ". ip:", val.ip, " database name:", val.database, " user id：", val.userid, " remark：", val.remark)
				i++
			}
		case "c":
			for _, v := range databaseip {
				if err := connMssqlPing(v); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(v.ip, v.database, v.userid, "链接成功")
				}
			}
		//case "dele":	支持删除
		default:
			break loop
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
func readINI(file, file_path, filename string) (map[string]connParamater, error) {
	var tp map[string]string
	tp = make(map[string]string)
	tp["DMS"] = "mssql"

	databaseip := make(map[string]connParamater)

	fileobj, err := os.Open(filepath.Join(file_path, filename))
	if err != nil {
		return databaseip, nil
	}
	defer fileobj.Close()

	readfileobj := bufio.NewReader(fileobj)
	readtxt, err := readfileobj.ReadString(byte('@'))
	if err != nil && err != io.EOF {
		return databaseip, err
	}

	splitstr := strings.Split(readtxt, "\n")
	for _, v := range splitstr {
		splitstr2 := strings.Split(v, "|")
		if len(splitstr2) >= 5 {
			if _, ok := tp[splitstr2[0]]; ok {
				databaseip[splitstr2[1]+"_"+splitstr2[2]] = connParamater{splitstr2[0], splitstr2[1], splitstr2[2], splitstr2[5], splitstr2[3], splitstr2[4]}
			}
		}
	}

	return databaseip, nil
}

func newConnMssqlPing(sc *bufio.Scanner) (connParamater, error) {
	var connParaSet connParamater

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

	connParaSet.tp = "DMS"

	if err := connMssqlPing(connParaSet); err != nil {
		return connParaSet, err
	} else {
		fmt.Println(connParaSet.ip, " ", connParaSet.database, "连接成功")

	}
	return connParaSet, nil
}

func connMssqlPing(connParaSet connParamater) error {
	connstr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", connParaSet.ip, connParaSet.userid, connParaSet.password, connParaSet.database)

	db, err := sql.Open("mssql", connstr)
	defer db.Close()

	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func saveConnStr(databaseip map[string]connParamater, file, file_path, filename string) error {
	var str string
	for _, connParaSet := range databaseip {
		str += fmt.Sprintf("%s|%s|%s|%s|%s|%s|\n", connParaSet.tp, connParaSet.ip, connParaSet.database, connParaSet.userid, connParaSet.password, connParaSet.remark)
	}

	if err := readwrite.WriteString(str, filepath.Join(file_path, filename)); err != nil {
		return err
	}
	return nil
}
