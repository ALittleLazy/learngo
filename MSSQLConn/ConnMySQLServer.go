package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)
func main() {
	sqlstr := "select stname from basestaff where stid = ?"

	db,err := sql.Open("mssql","server=127.0.0.1;user id =sa;password=ktwl;database=RHRDatabase;encrypt=disable")
	if err != nil{
		log.Println(err)
	}
	defer db.Close()

	err	= db.Ping()
	if err != nil {
		log.Println(err)
	}

	stmt,err := db.Prepare(sqlstr)
	defer stmt.Close()

	var name string
	id := 4

	err = stmt.QueryRow(id).Scan(&name)
	if err!= nil{
		log.Println(err)
	}
	fmt.Print("column name is ",name)

}
