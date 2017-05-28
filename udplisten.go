package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Host    string `short:"h" long:"host" default:"0.0.0.0" description:"Interface IP to bind to"`
	Port    uint16 `short:"p" long:"port" default:"1234" description:"UDP port to bind to"`
	File    string `short:"f" long:"file" default:"" description:"Append received data to"`
	Buffer  int    `short:"b" long:"buffer" default:"1500" description:"Max receive buffer size"`
	Newline bool   `short:"n" long:"newline" description:"Add newline to end of each message"`
	Quiet   bool   `short:"q" long:"quiet" description:"Suppress informational status on stderr"`
}

func handleClient(conn *net.UDPConn, fn string) {
	b := make([]byte, opts.Buffer)
	n, addr, e := conn.ReadFromUDP(b)
	if e != nil {
		log.Printf("Read from UDP failed, err: %v", e)
		return
	}
	if !opts.Quiet {
		log.Printf("Read from client(%v:%v), len: %v\n", addr.IP, addr.Port, n)
	}
	if opts.Newline && n < len(b) && b[n-1] != '\n' {
		b[n] = '\n'
		n += 1
	}
	if len(opts.File) != 0 {
		f, err := os.OpenFile(opts.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("Open file failed, err: %v", err)
			return
		}
		defer f.Close()
		if _, err = f.Write(b[:n]); err != nil {
			log.Printf("Write file failed, err: %v", err)
			return
		}
	} else { // Write to stdout
		if _, err := os.Stdout.Write(b[:n]); err != nil {
			log.Printf("Write stdout failed, err: %v", err)
			return
		}
	}
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("error: %v\n", err.Error())
		} else {
			os.Exit(0)
		}
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", opts.Host, opts.Port))
	if err != nil {
		log.Panic(err)
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
	}
	if !opts.Quiet {
		log.Printf("Starting udplistener at %v:%v", opts.Host, opts.Port)
	}
	for {
		handleClient(l, opts.File)
	}
}
