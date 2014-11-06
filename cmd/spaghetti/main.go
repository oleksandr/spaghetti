package main

import (
	"golang.org/x/net/websocket"
	"flag"
	"fmt"
	"github.com/oleksandr/spaghetti"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	bindAddr  = flag.String("bind", "0.0.0.0", "HTTP address to bind (default 0.0.0.0)")
	httpPort  = flag.Int("port", 3000, "HTTP port to listen (default 3000)")
	uplinkUrl = flag.String("uplink", "", "WS URL of another router to connect to")
)

func main() {
	flag.Parse()
	log.SetPrefix("[spaghetti] ")

	// Start signal catching routine for a propert exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func(ch chan os.Signal) {
		s := <-c
		log.Println("Got signal:", s)
		os.Exit(0)
	}(c)

	// Establish uplink connection is required
	if *uplinkUrl != "" {
		go spaghetti.DefaultHub.SetupUplink(*uplinkUrl)
	}

	// Map the URI to handlers
	http.Handle("/ws/pub", websocket.Handler(spaghetti.PubWebHandler))
	http.Handle("/ws/sub", websocket.Handler(spaghetti.SubWebHandler))
	http.Handle("/ws/pubsub", websocket.Handler(spaghetti.PubSubWebHandler))

	// Start the connections hub
	go spaghetti.DefaultHub.Start()

	// Listen & serve
	addr := fmt.Sprintf("%v:%v", *bindAddr, *httpPort)
	log.Printf("Listening %v", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe Error:", err)
	}
}
