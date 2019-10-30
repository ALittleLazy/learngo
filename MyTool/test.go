package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ttype struct {
	userid, password, remark string
}

var file, file_path, filename string

func main() {
	ip := make(map[string]map[string]ttype)
	databasename := make(map[string]ttype)
	databasename["test"] = ttype{"sa", "123", "shuoming"}
	ip["127.0.0.1"] = databasename
	fmt.Println(ip)

	v, ok := ip["127.0.0.1"]["test"]

	if ok == false {
		fmt.Println("not find", ok)
	} else {
		fmt.Println("find:", v)
	}

	file, _ := exec.LookPath(os.Args[0])
	file_path, _ := filepath.Abs(file)
	file_path, filename = filepath.Split(file_path)
	fmt.Println(file)
	fmt.Println(file_path)
	fmt.Println(filename)

	n := make(map[string]string)
	m := make(map[string]string)
	m["s"] = "s"
	m["s2"] = "s2"
	n["n"] = "n"
	n = m
	//delete(m, "s")
	fmt.Println(n, m)

}
