package mssql

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

type ConnParamater struct {
	Tp, Ip, Database, Remark, Userid, Password string
}

func NewConnMssqlPing(sc *bufio.Scanner) (ConnParamater, error) {
	var connParaSet ConnParamater

	fmt.Print("ip address:")
	sc.Scan()
	connParaSet.Ip = sc.Text()

	fmt.Print("database name:")
	sc.Scan()
	connParaSet.Database = sc.Text()

	fmt.Print("conn remark:")
	sc.Scan()
	connParaSet.Remark = sc.Text()

	fmt.Print("use id:")
	sc.Scan()
	connParaSet.Userid = sc.Text()

	fmt.Print("password:")
	sc.Scan()
	connParaSet.Password = sc.Text()

	connParaSet.Tp = "DMS"

	if err := ConnMssqlPing(connParaSet); err != nil {
		return connParaSet, err
	} else {
		fmt.Println(connParaSet.Ip, " ", connParaSet.Database, "连接成功")
	}
	return connParaSet, nil
}

func ConnMssqlPing(connParaSet ConnParamater) error {
	if connParaSet.Ip == "" {
		return errors.New("error:数据库地址为空")
	}

	connstr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", connParaSet.Ip, connParaSet.Userid, connParaSet.Password, connParaSet.Database)

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

func ConnMssqlExec(connParaSet ConnParamater, sqlstr string) ([]interface{}, error) {
	if connParaSet.Ip == "" {
		return nil, errors.New("error:数据库地址为空")
	}

	connstr := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", connParaSet.Ip, connParaSet.Userid, connParaSet.Password, connParaSet.Database)

	db, err := sql.Open("mssql", connstr)
	defer db.Close()

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	datarows, err := db.Query(sqlstr)
	defer datarows.Close()

	cols, _ := datarows.Columns()
	datainfo := make([]interface{}, 0, 30)
	rowinfo := make([]interface{}, len(cols)) //为什么make([]interface{}, len(cols))写为make([]interface{},0,len(cols))会报错？

	//为什么没有这部分会报错？
	row := make([]interface{}, len(cols))
	for idx := range rowinfo {
		row[idx] = new(interface{})
	}

	/*rowinfo = nil
	for _, val := range cols {
		rowinfo = append(rowinfo, val)
	}
	datainfo = append(datainfo, rowinfo)*/

	for datarows.Next() {

		err := datarows.Scan(row...)
		if err != nil {
			log.Println(err)
		}
		rowinfo = nil

		for _, val := range row {
			rowinfo = append(rowinfo, *(val.(*interface{}))) //val.(*interface{})这种写法不能理解是什么意思
		}
		datainfo = append(datainfo, rowinfo)
	}

	return datainfo, nil
}
