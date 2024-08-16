package main

import (
	"log"
	"net"
	"time"
)

const (
	poolSize = 10
)

type Job struct {
	conn net.Conn
}

func worker(id int, jobs <-chan Job) {
	log.Printf("Worker %d is processing a job.", id)

	for job := range jobs {
		handleConnection(job.conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	log.Println("Processing the request...")
	time.Sleep(5 * time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello World!\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal("Error getting listener")
		return
	}

	jobs := make(chan Job, poolSize)

	for w := 0; w < poolSize; w++ {
		go worker(w, jobs)
	}

	for {
		log.Println("Waiting for client to connect...")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error while getting connection object")
			return
		}
		jobs <- Job{conn: conn}
	}
}

/*
1. limiting the number of threads
=================================
const (
	maxThreads = 10
)
func handleConnection(conn net.Conn, sem chan struct{}) {
	// release semaphore slot when done
	defer func() {
		<-sem
		conn.Close()
	}()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	log.Println("Processing the request...")
	time.Sleep(5 * time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello World!\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":3000")

	if err != nil {
		log.Fatal("Error getting listener")
		return
	}

	sem := make(chan struct{}, maxThreads)

	for {
		log.Println("Waiting for client to connect...")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error while getting connection object")
			return
		}
		// take a slot
		sem <- struct{}{}
		go handleConnection(conn, sem)
	}
}


2. add thread pool to save on thread creation time
3. connection timeout
4. tcp backlog queue configuration
*/
