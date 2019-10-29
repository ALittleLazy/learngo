package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

type Match struct {
	MatchID          int64
}

func main() {
	var match Match
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable",
		"localhost", "sa", "ktwl", "test")

	conn, err := sql.Open("mssql", connString)
	conn2, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	defer conn2.Close()

	stmt, err := conn.Prepare("select top 1 * from t2")
	if err != nil {
		log.Printf("\nPrepare failed:%T %+v\n", err, err)

	}
	//defer	stmt.Close()
	row := stmt.QueryRow()
	//fmt.Print(row)
	err = row.Scan(&match.MatchID)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}
	fmt.Println(match)
	fmt.Println("SQL Drivers:",sql.Drivers())
}