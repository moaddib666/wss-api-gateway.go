package main

import (
	"MargayGateway/auth"
	"MargayGateway/backplane"
	"MargayGateway/constants"
	"MargayGateway/monitoring"
	"MargayGateway/registry"
	"github.com/gobwas/ws"
	"log"
	"net"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	monitoring.Init()
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:"+constants.DefaultPort)

	go monitoring.StartMetricsServer(constants.DefaultMetricsPort)
	if err != nil {
		log.Fatalf("can't start Service on 0.0.0.0:%s - %s", constants.DefaultPort, err)
	}
	ConnectionPool := registry.GetConnectionRegistry()
	EventBus := backplane.NewBus(ConnectionPool)

	for {
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
					connection = registry.CreateRegistryItem(client.Id(), conn)
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
	}
}
