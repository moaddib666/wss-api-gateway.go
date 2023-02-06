package main

import (
	"WSSFacade/auth"
	"WSSFacade/backplane"
	"WSSFacade/registry"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("can't start WSS on 0.0.0.0:8080")
	}
	ConnectionPool := registry.GetConnectionRegistry()
	EventBus := backplane.NewSampleBus(ConnectionPool)

	for {
		//time.Sleep(500)
		var connection *registry.Connection
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		u := ws.Upgrader{
			OnHeader: func(key, value []byte) error {
				if strings.ToLower(string(key)) != "authorization" {
					return nil
				}
				client := auth.NewClient(value)
				if client.IsAuthorised() {
					connection = registry.CreateRegistryItem(client.Id, conn)
					return nil
				}
				return ws.RejectConnectionError(
					ws.RejectionReason("bad authorization type"),
					ws.RejectionStatus(403),
				)
			},
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
			_ = conn.Close()
			continue
		}
		if connection == nil {
			log.Printf("Auth credetials was not provided")
			_ = conn.Close()
			continue
		}
		EventBus.ConnectClient(connection)
		//go handleConnection(connection)
	}
}

func handleConnection(connection *registry.Connection) {
	conn := connection.WebSocket
	defer conn.Close()

	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			// TODO handle error
			return
		}
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			// TODO handle error
			return
		}
		time.Sleep(5)
	}
}
