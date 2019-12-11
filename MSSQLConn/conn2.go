package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

func main() {
	var isdebug = true
	var server = "localhost"
	var port = 1433
	var user = "sa"
	var password = "ktwl"
	var database = "test"

	//连接字符串
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable", server, port, database, user, password)
	//connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable","localhost", "sa", "ktwl", "test")
	if isdebug {
		fmt.Println(connString)
	}
	//建立连接
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}
	defer conn.Close()

	//产生查询语句的Statement
	stmt, err := conn.Prepare(getMssqlJobStr())
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	//通过Statement执行查询
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}

	//建立一个列数组
	cols, err := rows.Columns()
	var colsdata = make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		colsdata[i] = new(interface{})
		//fmt.Print(cols[i])
		//fmt.Print("\t")
	}
	fmt.Println()

	//遍历每一行
	for rows.Next() {
		rows.Scan(colsdata...) //将查到的数据写入到这行中
		fmt.Println("colsdata:", colsdata)
		PrintRow(colsdata) //打印此行
	}
	defer rows.Close()
}

//打印一行记录，传入一个行的所有列信息
func PrintRow(colsdata []interface{}) {
	for _, val := range colsdata {
		switch v := (*(val.(*interface{}))).(type) {
		case nil:
			fmt.Print("NULL")
		case bool:
			if v {
				fmt.Print("True")
			} else {
				fmt.Print("False")
			}
		case []byte:
			fmt.Print(string(v))
		case time.Time:
			fmt.Print(v.Format("2016-01-02 15:05:05.999"))
		default:
			fmt.Print(v)
		}
		fmt.Print("\t")
	}
	fmt.Println()
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
		"SELECT c.name,current_execution_status = a.job_state,start_execution_date = (SELECT MAX(a2.start_execution_date) FROM msdb.dbo.sysjobactivity a2 WHERE a2.job_id = a.job_id AND a2.start_execution_date IS NOT NULL AND a2.stop_execution_date IS NULL AND a2.session_id = (SELECT MAX(a3.session_id) FROM msdb.dbo.syssessions a3)),last_run_outcome = (SELECT TOP 1 a2.last_run_outcome FROM msdb.dbo.sysjobservers a2 WHERE a2.job_id = a.job_id) FROM #help_job a JOIN msdb.dbo.sysjobs c ON c.job_id = a.job_id" + "\n" +
		"IF OBJECT_ID('tempdb..#help_job') IS NOT NULL DROP TABLE #help_job"
	//sqlstr = "select C1 from t1"
	return sqlstr
}
