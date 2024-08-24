package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	poolSize        = 10
	connectionLimit = 100
)

type Job struct {
	conn net.Conn
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Worker %d is processing a job.", id)

	for job := range jobs {
		log.Printf("Worker %d processing job: %+v\n", id, job.conn.RemoteAddr())
		handleConnection(job.conn)
	}
	log.Printf("Worker %d shutting down", id)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading from connection:", err)
		return
	}

	log.Println("Processing the request...")
	time.Sleep(1 * time.Second)

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello World!\r\n"))
	if err != nil {
		log.Println("Error writing to connection: ", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal("Error getting listener")
		return
	}
	defer listener.Close()

	jobs := make(chan Job, poolSize)
	var wg sync.WaitGroup

	for w := 0; w < poolSize; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// Handle shutdown signals
	shutdown := make(chan os.Signal, 1)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	connLimit := make(chan struct{}, connectionLimit) // Connection limiter

	go func() {
		<-shutdown
		log.Println("Shutdown signal received, stopping server...")
		listener.Close() // Close the listener to unblock Accept
		close(jobs)      // Close the job channel to stop workers
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-shutdown:
				// If the error is due to shutdown, exit the loop
				log.Println("Server is shutting down...")
				wg.Wait() // wait for all workers to finish
				log.Println("All workers have shut down gracefully.")
				return
			default:
				// If it's a different error, log and continue
				log.Println("Error accepting connection:", err)
			}
			continue
		}

		connLimit <- struct{}{} // Enforce connection limit

		go func(conn net.Conn) {
			defer func() { <-connLimit }() // Release connection slot
			jobs <- Job{conn: conn}
		}(conn)
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
