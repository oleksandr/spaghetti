package spaghetti

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"log"
	"strings"
)

// The communication hub between the active connections
type Hub struct {
	// Unique name for this hub
	ID string

	// All connections
	subscribers map[string]*Connection
	publishers  map[string]*Connection

	// Register requests from the connections
	Register chan *Connection

	// Unregister requests from connections
	Unregister chan *Connection

	// Uplink connection settings
	Uplink     *Connection
	UplinkType ConnectionType

	// Data handling channel
	Data chan *Message
}

// Start the hub and listen for incoming packets
func (self *Hub) Start() {
	log.Println("Hub.Start")
	for {
		select {
		case c := <-self.Register:
			log.Printf("Hub.Start: Register %v\n", c.ID)
			if c.Type == ConnectionTypePub || c.Type == ConnectionTypePubSub {
				self.publishers[c.ID] = c
			}
			if c.Type == ConnectionTypeSub || c.Type == ConnectionTypePubSub {
				self.subscribers[c.ID] = c
			}
			log.Printf("Publishers: %#v", self.publishers)
			log.Printf("Subscribers: %#v", self.subscribers)
		case c := <-self.Unregister:
			log.Printf("Hub.Start: Unregister %v\n", c.ID)
			delete(self.publishers, c.ID)
			delete(self.subscribers, c.ID)
		case m := <-self.Data:
			log.Println("Hub.Start: Data received")
			self.RouteMessage(m)
		}
	}
}

// Route a message to subscribers
func (self *Hub) RouteMessage(message *Message) {
	log.Println("Hub.RouteMessage")

	if len(self.subscribers) == 0 {
		return
	}

	if _, ok := self.publishers[message.ConnId]; !ok {
		log.Printf("Hub.RouteMessage: ERROR %v is not a publisher, ignoring...\n", message.ConnId)
	}

	//TODO: Transform message

	if self.Uplink != nil && self.Uplink.ID != message.ConnId && (self.UplinkType == ConnectionTypeSub || self.UplinkType == ConnectionTypePubSub) {
		select {
		case self.Uplink.Send <- message:
		default:
			close(self.Uplink.Send)
			go self.Uplink.WS.Close()
		}
	}

	for id, c := range self.subscribers {
		if c.ID == message.ConnId {
			// Don't sent itself
			continue
		}
		select {
		case c.Send <- message:
		default:
			delete(self.subscribers, id)
			delete(self.publishers, id)
			close(c.Send)
			go c.WS.Close()
		}
	}
}

// Establish uplink connection toa  given WS url
func (self *Hub) SetupUplink(uplinkUrl string) {
	log.Println("Hub.SetupUplink: ", uplinkUrl)

	origin := fmt.Sprintf("http://%v", self.ID)
	ws, err := websocket.Dial(uplinkUrl, "", origin)
	if err != nil {
		log.Printf("Hub.SetupUplink: ERROR: %v", err.Error())
		return
	}

	self.Uplink = NewConnection(ws)
	if strings.HasSuffix(uplinkUrl, "/pubsub") {
		self.UplinkType = ConnectionTypePubSub
	} else if strings.HasSuffix(uplinkUrl, "/pub") {
		self.UplinkType = ConnectionTypePub
	} else {
		self.UplinkType = ConnectionTypeSub
	}

	// Start read/write routines
	go self.Uplink.Reader()
	self.Uplink.Writer()
}

//
// The default Hub instance used all around this code
//
var DefaultHub = Hub{
	subscribers: make(map[string]*Connection),
	publishers:  make(map[string]*Connection),
	Register:    make(chan *Connection),
	Unregister:  make(chan *Connection),
	Data:        make(chan *Message),
}

//
// Init the package defaults
//
func init() {
	ID, _ := uuid.NewV4()
	DefaultHub.ID = ID.String()
}
