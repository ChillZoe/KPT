package gadget

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/sync/semaphore"
)

type FromTo struct {
	From int
	To   int
}

// var LivingAddress map[string]string

// LivingAddress := make([string]string[])

func GetIPList(ip string) (base string, start, end int, err error) {
	fromTo := strings.Split(ip, "-")
	ipStart := fromTo[0]
	err = fmt.Errorf("Invalid IP Range (eg. 1.1.1.1-3)\n")

	tIp := strings.Split(ipStart, ".")
	if len(tIp) != 4 {
		return
	}
	start, _ = strconv.Atoi(tIp[3])
	end = start
	if len(fromTo) == 2 {
		end, _ = strconv.Atoi(fromTo[1])
	}
	if end == 0 {
		return
	}
	base = fmt.Sprintf("%s.%s.%s", tIp[0], tIp[1], tIp[2])
	err = nil
	return
}

func GetPortList(s string) ([]FromTo, int) {
	res := make([]FromTo, 0)
	tot := 0

	for _, port := range strings.Split(s, ",") {
		from := 0
		to := 0
		fromTo := strings.Split(port, "-")
		from, _ = strconv.Atoi(fromTo[0])
		to = from
		if len(fromTo) == 2 {
			to, _ = strconv.Atoi(fromTo[1])
		}
		a := FromTo{
			From: from,
			To:   to,
		}
		res = append(res, a)
		tot += 1 + to - from
	}
	return res, tot
}

func Scan(ip string, port string, parallel int64, timeoutMS int) {
	portFromTo, _ := GetPortList(port)
	base, start, end, _ := GetIPList(ip)
	lock := semaphore.NewWeighted(parallel)
	timeout := time.Duration(timeoutMS) * time.Millisecond
	wg := sync.WaitGroup{}
	defer wg.Wait()
	// iterate ip in task list
	for ipExt := start; ipExt <= end; ipExt++ {
		ip := base + "." + fmt.Sprintf("%d", ipExt)
		// iterate port in task list
		for _, p := range portFromTo {
			// iterate port from A-B
			for port := p.From; port <= p.To; port++ {
				// lock down the context
				lock.Acquire(context.TODO(), 1)
				wg.Add(1)
				go func(port int, p FromTo) {
					defer lock.Release(1)
					defer wg.Done()

					if ScanPort(ip, port, timeout) {
						log.Info("Find ", ip, ":", port, " !\n")
						// LivingAddress[ip] = append(LivingAddress[ip], port)
					}

				}(port, p) // send all sync objects into args
			}
		}
	}
}

func ScanPort(ip string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			ScanPort(ip, port, timeout)
		}
		return false
	}

	_ = conn.Close()
	return true
}
