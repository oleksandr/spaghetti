Spaghetti
=========

A simple router (using Web-Sockets transport) to implement Pub/Sub messages delivery. This is not a real PubSub, not a Message Broker, not a Queue. There are no topics, exchanges, etc. This is a simple tool that helps to bind a bunch of devices and services in the office so they could communicate easily without too much network programming overhead.

## Installation

    go get github.com/oleksandr/spaghetti
    go install github.com/oleksandr/spaghetti/cmd/spaghetti

## Usage

Running the router:

    $ ./bin/spaghetti -bind=0.0.0.0 -port=3000

Connecting to another spaghetti router:

    $ ./bin/spaghetti -bind=0.0.0.0 -port=3100 -uplink=ws://localhost:3000/ws/pubsub

## Available endpoints

There are 3 endpoints for connecting to the router (either directly or using uplink):

 - Pub (ws://localhost:3000/ws/pub)
 - Sub (ws://localhost:3000/ws/sub)
 - PubSub (ws://localhost:3000/ws/pubsub)

As you can guess, Pub mode will allow the connected client/router only to send message to the network. Sub mode will only wait for message delivered from router. PubSub is for a bi-directional communication.

## Giving it a try

If you don't have Web-socket client or too lazy to program it just use the following URL [http://www.websocket.org/echo.html](http://www.websocket.org/echo.html) with echo client implemented to connect to your router.

## An example of a topology you can build with spaghetti

![Example](https://raw.github.com/oleksandr/spaghetti/master/example.png)

