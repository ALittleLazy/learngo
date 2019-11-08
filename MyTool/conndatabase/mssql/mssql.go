package mssql

import (
	"MyTool/aescbc"
	"MyTool/readwrite"
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type connParamater struct {
	tp, ip, database, remark, userid, password string
}

//var databaseip map[string] connParamater
const k string = "C70CB1D7A85944A08524065A4367392D"

var key, filekey string

func mainMenu() {
	fmt.Println()
	fmt.Println("n:新建连接；r:显示现有链接；c:测试连接 c [编号] ; d:删除连接 d [编号]|[all]")
}

func Menu() {
	var file, file_path, filename string
	file, _ = exec.LookPath(os.Args[0])
	file_path, _ = filepath.Abs(file)
	file_path, filename = filepath.Split(file_path)
	filename = strings.Replace(filename, filepath.Ext(filename), ".ini", -1)
	//filename = "___1go_build_CanDo_go.ini"

	menuno := make([]string, 1, 20)

	scanner := bufio.NewScanner(os.Stdin)

load:
	mainMenu()

	databaseip, err := readINI(file, file_path, filename)
	if err != nil {
		log.Println(err)
		filekey = k
		key = newKey(filename, filekey)
	}

	scanner.Scan()
	scanftxt := strings.Split(scanner.Text(), " ")
	switch strings.ToLower(scanftxt[0]) {
	case "n":
		newConnStr, err := newConnMssqlPing(scanner)
		if err != nil {
			log.Println(err)
		} else {
			databaseip[newConnStr.ip+"_"+newConnStr.database] = newConnStr
			if err := saveConnStr(databaseip, file, file_path, filename); err != nil {
				log.Println(err)
			}
		}
		goto load
	case "r":
		i := 0
		if len(menuno) < len(databaseip) {
			menuno = append(menuno, make([]string, len(databaseip)-len(menuno)+10)...)
		}
		for no, val := range databaseip {
			menuno[i] = no
			fmt.Println(i+1, ". ip:", val.ip, " database name:", val.database, " user id：", val.userid, " remark：", val.remark)
			i++
		}
		goto load
	case "c":
		if len(scanftxt) > 1 {
			for _, v := range scanftxt {
				if len(v) > 0 {
					if ok, _ := regexp.MatchString("\\d", v); ok {
						j, _ := strconv.ParseInt(v, 0, 0)
						j = j - 1
						if j >= 0 && j <= int64(len(menuno)-1) {
							if err := connMssqlPing(databaseip[menuno[j]]); err != nil {
								fmt.Println(err)
							} else {
								fmt.Println(databaseip[menuno[j]].ip, databaseip[menuno[j]].database, databaseip[menuno[j]].userid, "链接成功")
							}
						}
						goto load
					}
				}
			}
		}
		for _, v := range databaseip {
			if err := connMssqlPing(v); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(v.ip, v.database, v.userid, "链接成功")
			}
		}
		goto load
	case "d":
		if len(scanftxt) > 1 {
			for _, v := range scanftxt {
				if len(v) > 0 {
					if ok, _ := regexp.MatchString("\\d", v); ok {
						j, _ := strconv.ParseInt(v, 0, 0)
						j = j - 1
						if j >= 0 && j <= int64(len(menuno)-1) {
							msg := fmt.Sprintf("%s %s %s", databaseip[menuno[j]].ip, databaseip[menuno[j]].database, databaseip[menuno[j]].userid)
							delete(databaseip, menuno[j])
							if _, ok := databaseip[menuno[j]]; ok {
								fmt.Println(msg, "删除成功")
							}
						}
					} else if "all" == strings.ToLower(v) {
						databaseip = nil
						fmt.Println("列表已清空")
					} else {
						continue
					}
					if err := saveConnStr(databaseip, file, file_path, filename); err != nil {
						log.Println(err)
					}
					break
				}
			}
		}
		goto load
	default:
		break
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

	//fileobj, err := os.Open(filepath.Join(file_path, filename))
	//if err != nil {
	//	return databaseip, nil
	//}
	//defer fileobj.Close()

	//readfileobj := bufio.NewReader(fileobj)
	//readtxt, err := readfileobj.ReadString(byte('@'))
	//if err != nil && err != io.EOF {
	//	return databaseip, err
	//}
	readtxt, err := readwrite.ReadString(filepath.Join(file_path, filename))
	if err != nil && err != io.EOF {
		return databaseip, err
	}

	splitstr := strings.Split(readtxt, "\n")
	for _, v := range splitstr {
		if len(v) >= 4 {
			if v[:4] != "KEY|" {
				continue
			}
		}
		splitstr2 := strings.Split(v, "|")
		if len(splitstr2) >= 2 {
			if splitstr2[0] == "KEY" {
				if len(splitstr2[1]) < 32 {
					filekey = k
				} else {
					filekey = splitstr2[1][:32]
				}

				key = newKey(filename, filekey)
				break
			}
		}
	}

	for _, v := range splitstr {

		if len(v) >= 4 {
			if v[:4] == "KEY|" {
				continue
			}
		} else {
			continue
		}

		v = aescbc.AesDecrypt(v, key)
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
	if connParaSet.ip == "" {
		return errors.New("error:数据库地址为空")
	}

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
	var str string = "KEY|" + key + "\n"

	for _, connParaSet := range databaseip {
		str += aescbc.AesEncrypt(fmt.Sprintf("%s|%s|%s|%s|%s|%s|", connParaSet.tp, connParaSet.ip, connParaSet.database, connParaSet.userid, connParaSet.password, connParaSet.remark), newKey(filename, key)) + "\n"
	}

	if err := readwrite.WriteString(str, filepath.Join(file_path, filename)); err != nil {
		return err
	}
	return nil
}
func newKey(txt string, key string) string {
	str := aescbc.AesEncrypt(txt, key)
	if len(str) >= 32 {
		str = str[len(str)-32:]
	} else {
		str = string(aescbc.PKCS7Padding([]byte(str), 32))
	}
	return str
}
