package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Input_appName string
var Input_Hour int
var Input_Minute int
var Input_Second int
var Input_Query string

var wg sync.WaitGroup

func Init() {
	flag.StringVar(&Input_appName, "app", "auth-dev", "Input app name")
	flag.IntVar(&Input_Hour, "h", 0, "Input Hour")
	flag.IntVar(&Input_Minute, "m", 0, "Input Minute")
	flag.IntVar(&Input_Second, "s", 0, "Input Second")
	flag.StringVar(&Input_Query, "q", "", "Input Query")
}

func main() {
	Init()
	flag.Parse()

	timeToGetLogs()
	//doGetLogs(1)
	//doGetLogs(2)
	//doGetLogs(3)
	//doGetAllLogs()
	wg.Wait()
}

func timeToGetLogs() {
	//oNum := 2
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), Input_Hour, Input_Minute, Input_Second, 0, now.Location())
		if next.Before(now) {
			next = now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), Input_Hour, Input_Minute, Input_Second, 0, next.Location())
		}

		t := time.NewTimer(next.Sub(now))
		<-t.C
		wg.Add(1)
		go doGetAllLogs()
		//doGetLogs(1)
		//t2 := time.NewTicker(time.Minute * 15)
		//for {
		//	select {
		//	case <-t2.C:
		//		doGetLogs(oNum)
		//		oNum++
		//	}
		//}
		break
	}
}

func doGetAllLogs() {
	defer wg.Done()

	//ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()
	//cmd := exec.CommandContext(ctx, "C:/Program Files/Cloud Foundry/cf.exe", "logs", Input_appName)
	cmd := exec.Command("C:/Program Files/Cloud Foundry/cf.exe", "logs", Input_appName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("cf開始にエラーが発生した：　%s\n", err.Error())
	}
	bf := bufio.NewScanner(stdout)
	fName := "./" + Input_appName + ".log"
	file, err := os.OpenFile(fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		fmt.Printf("ファイル作成にエラーが発生した：　%s\n", err.Error())
	}
	wr := bufio.NewWriter(file)
	for bf.Scan() {
		if line := bf.Text(); strings.Contains(line, Input_Query) {
			if _, err := wr.WriteString(line+"\n"); err != nil {
				fmt.Printf("書き込みにエラーが発生した：　%s\n", err.Error())
			}
			wr.Flush()
		}
	}
	//bs, _ := ioutil.ReadAll(stdout)
	//fName := "./"+Input_appName+".log"
	//if fileExists(fName) {
	//	fName = "./"+Input_appName+".log"
	//}
	//if err := ioutil.WriteFile(fName, bs, 0644); err != nil {
	//	fmt.Printf("書き込みにエラーが発生した：　%s\n", err.Error())
	//}

}

func doGetLogs(numb int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "C:/Program Files/Cloud Foundry/cf.exe", "logs", Input_appName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("cf開始にエラーが発生した：　%s\n", err.Error())
	}
	bs, _ := ioutil.ReadAll(stdout)
	fName := "./" + Input_appName + ".log"
	if fileExists(fName) {
		fName = "./" + Input_appName + strconv.Itoa(numb) + ".log"
	}
	if err := ioutil.WriteFile(fName, bs, 0644); err != nil {
		fmt.Printf("書き込みにエラーが発生した：　%s\n", err.Error())
	}

}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
