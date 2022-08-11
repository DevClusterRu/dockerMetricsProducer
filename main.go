package main

import (
	"fmt"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"log"
	"math"
	"net/rpc"
	"os"
	"time"
)

type System struct {
	Node   string
	Cpu    float64
	Memory float64
}

func main() {

	ip, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	for {
		sys := System{
			Node: ip,
		}

		memory, err := memory.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		sys.Memory = math.Round(float64(memory.Used/1024) / float64(memory.Total/1024) * 100)
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
		sys.Cpu = 100 - float64(after.Idle-before.Idle)/total*100

		client, err := rpc.DialHTTP("tcp", "3.134.16.137:4480")
		if err != nil {
			log.Panicf("Error in dialing. %s", err.Error())
		}
		defer client.Close()
		var result string

		err = client.Call("Jobs.PushMetric", sys, &result)
		if err != nil {
			log.Println("Error in push. %s", err.Error())
		}

		time.Sleep(5 * time.Second)

	}

}
