package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	interval int
	failinit int
)

func main() {
	interval, _ = strconv.Atoi(os.Getenv("INTERVAL"))
	failinit, _ = strconv.Atoi(os.Getenv("FAILINIT"))
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	LOG_FILE := path + "/log/" + time.Now().Format("01-02-2006") + ".log"
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	host1 := os.Getenv("HOST1")
	host2 := os.Getenv("HOST2")
	host3 := os.Getenv("HOST3")
	port := "9000"
	var faill int = 0

	for {
		var check1 bool = checkport(host1, port)
		var check2 bool = checkport(host2, port)
		var check3 bool = checkport(host3, port)
		currentTime := time.Now().Format("01-02-2006 15:04:05")
		// if check1 == true || check2 == true || check3 == true {
		if check1 == true || check2 == true || check3 == true {
			fmt.Println("[" + currentTime + "] [INFO] No need Failover")
			log.Println("[INFO] No need Failover")
			faill = 0
		} else {
			faill++
			fmt.Println("["+currentTime+"] [WARN] Need Failover ,TIMES : ", faill)
			log.Println("[WARN] Need Failover ,", faill)
			if faill >= failinit {
				fmt.Println("[" + currentTime + "] [FAIL] TIME OUT GO Failover....")
				log.Println("[FAIL] TIME OUT GO Failover....")
				promote()
				break
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func promote() {
	run, err := exec.Command("sudo", "systemctl", "start", "sonarqube").Output()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(run))
}

func checkport(url_host string, port string) bool {
	target := net.JoinHostPort(url_host, port)
	timeout := 3 * time.Second
	currentTime := time.Now().Format("01-02-2006 15:04:05")
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		fmt.Println("[" + currentTime + "] Sonarqube on " + url_host + " Down")
		return false
	} else {
		if conn != nil {
			fmt.Println("[" + currentTime + "] Sonarqube on " + url_host + " UP")
			_ = conn.Close()
			return true
		} else {
			fmt.Println("[" + currentTime + "] Sonarqube on " + url_host + " Down")
			return false
		}
	}
}
