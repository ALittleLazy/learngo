package readwrite

import (
	"bufio"
	"io"
	"os"
)

func WriteString(wstr string, filename string) error {
	var fileobj *os.File
	var err error

	fileobj, err = os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer fileobj.Close()

	if err != nil {
		//log.Println(err)
		return err
	}

	writeobj := bufio.NewWriter(fileobj)
	_, err = writeobj.WriteString(wstr)
	if err != nil {
		//log.Println(err)
		return err
	}

	if err = writeobj.Flush(); err != nil {
		return err
	}

	return nil
}
func ReadString(filename string) (str string, err error) {
	fileobj, err := os.Open(filename)
	if err != nil {
		return str, err
	}
	defer fileobj.Close()

	readfileobj := bufio.NewReader(fileobj)
	readtxt, err := readfileobj.ReadString(byte('@'))
	if err != nil && err != io.EOF {
		return str, err
	}

	str = readtxt
	return str, nil
}
