package main

import "fmt"

type ttype struct {
	userid, password, remark string
}

func main() {
	ip := make(map[string]map[string]ttype)
	databasename := make(map[string]ttype)
	databasename["test"] = ttype{"sa", "123", "shuoming"}
	ip["127.0.0.1"] = databasename
	fmt.Println(ip)

	v, ok := ip["127.0.0.1"]["test"]

	if ok == false {
		fmt.Print("not find", ok)
	} else {
		fmt.Print("find:", v)
	}

}
