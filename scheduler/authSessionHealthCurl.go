package main

import (
	"flag"
	"fmt"
	"os/exec"
	"time"
)
var Hour_Sesson int
var Minute_Sesson int
var Second_Sesson int
func main() {
	flag.IntVar(&Hour_Sesson,"h",0,"Input Hour")
	flag.IntVar(&Minute_Sesson,"m",0,"Input Minute")
	flag.IntVar(&Second_Sesson,"s",0,"Input Second")
	flag.Parse()

	timeToSessionCurlExecute()
	//doSessionCurlPeriodically()
}

func timeToSessionCurlExecute() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), Hour_Sesson, Minute_Sesson, Second_Sesson, 0, now.Location())
		if next.Before(now) {
			next = now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), Hour_Sesson, Minute_Sesson, Second_Sesson, 0, next.Location())
		}

		t := time.NewTimer(next.Sub(now))
		<-t.C
		doSessionCurlPeriodically()
		break
	}
}

func doSessionCurlPeriodically() {
	cmd := exec.Command("./sessionCurl.sh")

	if err := cmd.Run(); err != nil {
		fmt.Printf("エラーが発生した：　%s\n", err.Error())
	}
}
