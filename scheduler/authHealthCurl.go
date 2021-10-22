package main

import (
	"context"
	"flag"
	"fmt"
	"os/exec"
	"time"
)

var Hour_Curl int
var Minute_Curl int
var Second_Curl int

func main() {
	flag.IntVar(&Hour_Curl, "h", 0, "Input Hour")
	flag.IntVar(&Minute_Curl, "m", 0, "Input Minute")
	flag.IntVar(&Second_Curl, "s", 0, "Input Second")
	flag.Parse()

	timeToCurlExecute()
	//doCurlPeriodically()
}

func timeToCurlExecute() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), Hour_Curl, Minute_Curl, Second_Curl, 0, now.Location())
		if next.Before(now) {
			next = now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), Hour_Curl, Minute_Curl, Second_Curl, 0, next.Location())
		}

		t := time.NewTimer(next.Sub(now))
		<-t.C
		doCurlPeriodically()
		t2 := time.NewTicker(time.Minute * 15)
		for {
			select {
			case <-t2.C:
				doCurlPeriodically()
			}
		}
		//t3 := time.NewTimer(time.Minute * 30)
		//<-t3.C
		//doCurlPeriodically()
		//break
	}
}

func doCurlPeriodically() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx,"./authCurl.sh")

	if err := cmd.Run(); err != nil {
		fmt.Printf("エラーが発生した：　%s\n", err.Error())
	}
}
