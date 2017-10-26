package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

type DataResponse struct {
	Hostname    string      `json:"hostname,omitempty"`
	IP          []string    `json:"ip,omitempty"`
	Headers     http.Header `json:"headers,omitempty"`
	Environment []string    `json:"environment,omitempty"`
}

var port string

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()

	http.HandleFunc("/", whoamI)
	http.HandleFunc("/api", api)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fetchData(req *http.Request) DataResponse {
	hostname, _ := os.Hostname()
	data := DataResponse{
		hostname,
		[]string{},
		req.Header,
		os.Environ(),
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			data.IP = append(data.IP, ip.String())
		}
	}

	return data
}

func printBinary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		printBinary(p)
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func whoamI(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	queryParams := u.Query()

	wait := queryParams.Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	data := fetchData(req)
	fmt.Fprintln(w, "Hostname:", data.Hostname)

	for _, ip := range data.IP {
		fmt.Fprintln(w, "IP:", ip)
	}

	for _, env := range data.Environment {
		fmt.Fprintln(w, "ENV:", env)
	}

	req.Write(w)
}

func api(w http.ResponseWriter, req *http.Request) {
	data := fetchData(req)
	json.NewEncoder(w).Encode(data)
}

type healthState struct {
	StatusCode int
}

var currentHealthState = healthState{200}
var mutexHealthState = &sync.RWMutex{}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var statusCode int
		err := json.NewDecoder(req.Body).Decode(&statusCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			fmt.Printf("Update health check status code [%d]\n", statusCode)
			mutexHealthState.Lock()
			defer mutexHealthState.Unlock()
			currentHealthState.StatusCode = statusCode
		}
	} else {
		mutexHealthState.RLock()
		defer mutexHealthState.RUnlock()
		w.WriteHeader(currentHealthState.StatusCode)
	}
}
