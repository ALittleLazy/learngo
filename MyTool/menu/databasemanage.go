package menu

import (
	"MyTool/aescbc"
	"MyTool/conndatabase/mssql"
	"MyTool/readwrite"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const k string = "C70CB1D7A85944A08524065A4367392D"

var key, filekey string

func mainMenu() {
	fmt.Println()
	fmt.Println("n:新建连接；r:显示现有链接；c:测试连接 c [编号][a:作业检查(最后一次)] ; d:删除连接 d [编号]|[all]")
}

func Menu() {
	var file, file_path, filename string
	file, _ = exec.LookPath(os.Args[0])
	file_path, _ = filepath.Abs(file)
	file_path, filename = filepath.Split(file_path)
	filename = strings.Replace(filename, filepath.Ext(filename), ".ini", -1)
	//filename = "___1go_build_CanDo_go.ini"

	menuno := make([]string, 1, 20)
	//var jobsinfo = make([]interface{}, 0, 100)

	scanner := bufio.NewScanner(os.Stdin)

load:
	mainMenu()

	databaseip, err := readINI(file, file_path, filename)
	if err != nil {
		log.Println(err)
		filekey = k
		key = NewKey(k, filekey)
	}

	scanner.Scan()
	scanftxt := strings.Split(scanner.Text(), " ")
	switch strings.ToLower(scanftxt[0]) {
	case "n": //新建
		newConnStr, err := mssql.NewConnMssqlPing(scanner)
		if err != nil {
			log.Println(err)
		} else {
			databaseip[newConnStr.Ip+"_"+newConnStr.Database] = newConnStr
			if err := saveConnStr(databaseip, file, file_path, filename); err != nil {
				log.Println(err)
			}
		}
		goto load
	case "r": //加载配置文件
		i := 0
		if len(menuno) < len(databaseip) {
			menuno = append(menuno, make([]string, len(databaseip)-len(menuno)+10)...)
		}
		for no, val := range databaseip {
			menuno[i] = no
			fmt.Println(i+1, ". ip:", val.Ip, " database name:", val.Database, " user id：", val.Userid, " remark：", val.Remark)
			i++
		}
		goto load
	case "c": //测试链接
		if len(scanftxt) > 1 {
			var j int64
			var item string

			for idx, v := range scanftxt {
				if idx > 0 && len(v) > 0 {
					if ok, _ := regexp.MatchString("\\d", v); ok {
						j, _ = strconv.ParseInt(v, 0, 0)
					}
					if ok, _ := regexp.MatchString("\\D", v); ok {
						item = v
					}
				}
			}

			if j > 0 && j <= int64(len(menuno)-1) && item == "" {
				j--
				if err := mssql.ConnMssqlPing(databaseip[menuno[j]]); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(databaseip[menuno[j]].Ip, databaseip[menuno[j]].Database, databaseip[menuno[j]].Userid, "链接成功")
				}
			} else if j > 0 && j <= int64(len(menuno)-1) && item == "a" {
				//fmt.Println(getMssqlJobStr())
				j--
				if datainfo, err := mssql.ConnMssqlExec(databaseip[menuno[j]], getMssqlJobStr()); err != nil {
					fmt.Println(err)
				} else {

					for _, val := range datainfo {
						/*取传值中的inteface{}各项值
						for idx:=0;idx < reflect.ValueOf(val).Len();idx++{
							fmt.Println("idx value:",reflect.ValueOf(val).Index(idx))
						}*/
						/*v := []interface{}{databaseip[menuno[j]].Ip}
						val := append(v,val)
						jobsinfo = append(jobsinfo,val)*/
						fmt.Println(databaseip[menuno[j]].Ip, val)
					}
					//fmt.Println(datainfo[0][2])
				}
			} else {
				for _, v := range databaseip {
					if datainfo, err := mssql.ConnMssqlExec(v, getMssqlJobStr()); err != nil {
						fmt.Println(err)
					} else {
						for _, val := range datainfo {
							/*取传值中的inteface{}各项值
							for idx:=0;idx < reflect.ValueOf(val).Len();idx++{
								fmt.Println("idx value:",reflect.ValueOf(val).Index(idx))
							}*/
							/*v := []interface{}{databaseip[menuno[j]].Ip}
							val := append(v,val)
							jobsinfo = append(jobsinfo,val)*/
							fmt.Println(v.Ip, val)
						}
						//fmt.Println(datainfo[0][2])
					}
				}
			}
			goto load
		}
		for _, v := range databaseip {
			if err := mssql.ConnMssqlPing(v); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(v.Ip, v.Database, v.Userid, "链接成功")
			}
		}
		goto load
	case "d": //删除
		if len(scanftxt) > 1 {
			for _, v := range scanftxt {
				if len(v) > 0 {
					if ok, _ := regexp.MatchString("\\d", v); ok {
						j, _ := strconv.ParseInt(v, 0, 0)
						j = j - 1
						if j >= 0 && j <= int64(len(menuno)-1) {
							msg := fmt.Sprintf("%s %s %s", databaseip[menuno[j]].Ip, databaseip[menuno[j]].Database, databaseip[menuno[j]].Userid)
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

func readINI(file, file_path, filename string) (map[string]mssql.ConnParamater, error) {
	var tp map[string]string
	tp = make(map[string]string)
	tp["DMS"] = "mssql"

	databaseip := make(map[string]mssql.ConnParamater)

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

				key = NewKey(k, filekey)
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
				databaseip[splitstr2[1]+"_"+splitstr2[2]] = mssql.ConnParamater{splitstr2[0], splitstr2[1], splitstr2[2], splitstr2[5], splitstr2[3], splitstr2[4]}
			}
		}
	}
	return databaseip, nil
}

func getMssqlJobStr() string {
	var sqlstr string
	sqlstr = "create table #help_job(" + "\n" +
		"job_id uniqueidentifier not null," + "\n" +
		"last_run_date int not null," + "\n" +
		"last_run_time int not null," + "\n" +
		"next_run_date int not null," + "\n" +
		"next_run_time int not null," + "\n" +
		"next_run_schedule_id int not null," + "\n" +
		"requested_to_run int not null, " + "\n" +
		"request_source int not null," + "\n" +
		"request_source_id sysname null," + "\n" +
		"running int not null, " + "\n" +
		"current_step int not null," + "\n" +
		"current_retry_attempt int not null," + "\n" +
		"job_state int not null" + "\n" +
		")" + "\n" +
		"insert into #help_job execute master.dbo.xp_sqlagent_enum_jobs 1, 'sa'" + "\n" +
		"SELECT CASE e.enabled WHEN 1 THEN '启用' ELSE '未启用' END,e.name,CASE e.current_execution_status WHEN 1 THEN '执行中' WHEN 3 THEN '重试间隔' WHEN 4 THEN '当前空闲' END,e.start_execution_date,CASE e.last_run_outcome WHEN 0 THEN '上次执行失败' WHEN 1 THEN '上次执行成功' WHEN 3 THEN '上次执行取消' WHEN 5 THEN '未知' END FROM (SELECT c.name,current_execution_status = a.job_state,start_execution_date = (SELECT MAX(a2.start_execution_date) FROM msdb.dbo.sysjobactivity a2 WHERE a2.job_id = a.job_id AND a2.start_execution_date IS NOT NULL AND a2.stop_execution_date IS NULL AND a2.session_id = (SELECT MAX(a3.session_id) FROM msdb.dbo.syssessions a3)),last_run_outcome = (SELECT TOP 1 a2.last_run_outcome FROM msdb.dbo.sysjobservers a2 WHERE a2.job_id = a.job_id),c.enabled FROM #help_job a JOIN msdb.dbo.sysjobs c ON c.job_id = a.job_id WHERE a.job_state IN (3,4)) e WHERE e.last_run_outcome IN (0,3) AND e.enabled = 1" + "\n" +
		"IF OBJECT_ID('tempdb..#help_job') IS NOT NULL DROP TABLE #help_job"
	//sqlstr = "select C1 from t1"
	return sqlstr
}

func saveConnStr(databaseip map[string]mssql.ConnParamater, file, file_path, filename string) error {
	var str string = "KEY|" + key + "\n"

	for _, connParaSet := range databaseip {
		str += aescbc.AesEncrypt(fmt.Sprintf("%s|%s|%s|%s|%s|%s|", connParaSet.Tp, connParaSet.Ip, connParaSet.Database, connParaSet.Userid, connParaSet.Password, connParaSet.Remark), NewKey(k, key)) + "\n"
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
