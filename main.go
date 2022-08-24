package main

import (
	"fmt"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"log"
	"math"
	"net/rpc"
	"os"
	"strings"
	"time"
)

type System struct {
	Node        string
	MetricName  string
	MetricValue float64
}

func main() {

	ip, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	rpcAddr := "172.31.47.31:4480"

	if strings.Contains(ip, "dcdell") {
		rpcAddr = "3.134.16.137:4480"
	}

	client, err := rpc.DialHTTP("tcp", rpcAddr)
	if err != nil {
		log.Panicf("Error in dialing. %s", err.Error())
	}

	log.Println("Started....")

	defer client.Close()

	for {

		mPackage := []System{}

		memory, err := memory.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		mem := math.Round(float64(memory.Used/1024) / float64(memory.Total/1024) * 100)

		mPackage = append(mPackage, System{
			Node:        ip,
			MetricName:  "memory_available_percently",
			MetricValue: mem,
		})

		before, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		time.Sleep(time.Duration(1) * time.Second)
		after, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		total := float64(after.Total - before.Total)

		cpu := 100 - float64(after.Idle-before.Idle)/total*100

		mPackage = append(mPackage, System{
			Node:        ip,
			MetricName:  "cpu_available_percently",
			MetricValue: cpu,
		})

		var result string

		err = client.Call("Jobs.PushMetric", mPackage, &result)
		if err != nil {
			log.Fatal("Error in push. %s", err.Error())
		}

		time.Sleep(5 * time.Second)

	}

}
