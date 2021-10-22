package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func main() {
	//srcName := os.Args[1]
	//dstName := os.Args[2]

	//copyFiles(srcName)
	copyFiles("SSO.xlsx")

	wg.Wait()
}

func dirWalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, filepath.Join(dir, file.Name()))
			//paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
		}

	}

	return paths
}

func copyFiles(srcName string) {
	bs, err := ioutil.ReadFile(srcName)
	if err != nil {
		fmt.Printf("エクセルファイル読み込みにエラーが発生した：　%s\n", err.Error())
	}

	for _, fd := range dirWalk(".") {
		wg.Add(1)
		go worker(fd, srcName, bs)
	}
}

func worker(fd string, srcName string, bs []byte) {
	defer wg.Done()

	ioutil.WriteFile(fd+"/"+srcName, bs, 0644)

	file, err := os.OpenFile(fd+"/throughput.txt", os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Printf("ファイル読み込みにエラーが発生した：　%s\n", err.Error())
	}
	f, err := excelize.OpenFile(fd + "/" + srcName)
	if err != nil {
		fmt.Println(err)
		return
	}
	bf := bufio.NewScanner(file)
	num := 2
	for bf.Scan() {
		line := bf.Text()
		cell := strings.Split(line, "\t")
		if len(cell) != 5 {
			continue
		}
		f.SetCellValue("LogData", "A"+strconv.Itoa(num), cell[0])
		f.SetCellValue("LogData", "B"+strconv.Itoa(num), cell[1])
		f.SetCellValue("LogData", "C"+strconv.Itoa(num), cell[2])
		f.SetCellValue("LogData", "D"+strconv.Itoa(num), cell[3])
		num++
	}
	f.UpdateLinkedValue()
	// Save the xlsx file with the origin path.
	if err = f.Save(); err != nil {
		fmt.Println(err)
	}
}
