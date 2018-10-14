package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var address string

const addrDef string = "golang.org"
const addrDesc string = "hostname or IP address"

var port int

const portDef int = 80
const portDesc = "port"

var timeout int

const timeoutDef int = 3
const timeoutDesc string = "how long to wait before timing out and terminating the smell"

var httpReq, httpsReq bool

const httpDef bool = false
const httpDesc string = "attempt an http GET request"
const httpsDef bool = false
const httpsDesc string = "attempt an https GET request"

func main() {
	fmt.Println("starting smeller")
	flag.StringVar(&address, "address", addrDef, addrDesc)
	flag.StringVar(&address, "a", addrDef, addrDesc+" (shorthand)")
	flag.IntVar(&port, "port", portDef, portDesc)
	flag.IntVar(&port, "p", portDef, portDesc+" (shorthand)")
	flag.IntVar(&timeout, "timeout", timeoutDef, timeoutDesc)
	flag.IntVar(&timeout, "t", timeoutDef, timeoutDesc+" (shorthand)")
	flag.BoolVar(&httpReq, "http", httpDef, httpDesc)
	flag.BoolVar(&httpsReq, "https", httpsDef, httpsDesc)
	flag.Parse()

	ip, err := net.LookupIP(address)
	fmt.Printf("ip addresses: %s\n", ip)
	cname, err := net.LookupCNAME(address)
	fmt.Printf("cname: %s\n", cname)
	dnsNames, err := net.LookupAddr(address)
	fmt.Printf("dns names: %s\n", dnsNames)
	isAlive, elapsedTimeMessage, err := tcpAlive()
	if err == nil {
		fmt.Printf("tcp port alive:\t %t\n", isAlive)
		fmt.Println(elapsedTimeMessage)
	} else {
		fmt.Println("no response from port")
		os.Exit(1)
	}
	if httpReq {
		respStatus, respTime := httpGet("http://")
		fmt.Printf("http status: %s\n", respStatus)
		fmt.Println(respTime)
	}
	if httpsReq {
		respStatus, respTime := httpGet("https://")
		fmt.Printf("https status: %s\n", respStatus)
		fmt.Println(respTime)
	}
}

func tcpAlive() (bool, string, error) {
	start := time.Now()
	one := []byte{}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(address, strconv.Itoa(port)), time.Second*3)
	if err == nil {
		conn.SetReadDeadline(time.Now())
		if _, err := conn.Read(one); err == io.EOF {
			conn.Close()
			conn = nil
			return false, elapsedTimeMessage(start), err
		}
		return true, elapsedTimeMessage(start), nil
	}
	return false, elapsedTimeMessage(start), err
}

func elapsedTimeMessage(start time.Time) string {
	elapsed := time.Since(start)
	return fmt.Sprintf("took %s", elapsed)
}

func httpGet(protocol string) (string, string) {
	start := time.Now()
	fullAddr := protocol + net.JoinHostPort(address, strconv.Itoa(port))
	resp, err := http.Get(fullAddr)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	respTime := elapsedTimeMessage(start)
	//body, err := ioutil.ReadAll(resp.Body)
	//fmt.Print(string(body[:]))

	return resp.Status, respTime
}
