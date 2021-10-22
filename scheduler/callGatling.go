package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

var Num_Call int
var Interval_Call int
var Hour_Call int
var Minute_Call int
var Second_Call int

func main() {
	flag.IntVar(&Num_Call, "num", 1, "Input Num")
	flag.IntVar(&Interval_Call, "i", 15, "Input Interval(Minutes)")
	flag.IntVar(&Hour_Call, "h", 0, "Input Hour")
	flag.IntVar(&Minute_Call, "m", 0, "Input Minute")
	flag.IntVar(&Second_Call, "s", 0, "Input Second")
	flag.Parse()

	timeToExecute()
	//doPeriodically()
}

func timeToExecute() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), Hour_Call, Minute_Call, Second_Call, 0, now.Location())
		if next.Before(now) {
			next = now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), Hour_Call, Minute_Call, Second_Call, 0, next.Location())
		}

		t := time.NewTimer(next.Sub(now))
		<-t.C
		doPeriodically()

		t2 := time.NewTicker(time.Minute * time.Duration(Interval_Call))
		for {
			select {
			case <-t2.C:
				doPeriodically()
			}
		}

		//t3 := time.NewTimer(time.Minute * 30)
		//<-t3.C
		//doPeriodically()
		//break
	}
}

func doPeriodically() {
	cmd := exec.Command("./gatling.sh")
	stdin, _ := cmd.StdinPipe()

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(strconv.Itoa(Num_Call) + "\n"))
		stdin.Write([]byte("\n"))
		stdin.Write([]byte("\n"))
	}()

	if err := cmd.Start(); err != nil {
		fmt.Printf("エラーが発生した：　%s\n", err.Error())
	}
}
