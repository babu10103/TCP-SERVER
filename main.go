package main

import (
	"log"
	"net"
	"time"
)

func do(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Error")
	}
	time.Sleep(1 * time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello World!\r\n"))

	conn.Close()

}

func main() {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal("Error getting listener")
		return
	}

	for {
		log.Println("Waiting for client to connect...")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error while getting connection object")
			return
		}
		log.Println("Connected successfully.")
		go do(conn)
	}
}

/*
1. limiting the number of threads
2. add thread pool to save on thread creation time
3. connection timeout
4. tcp backlog queue configuration
*/
