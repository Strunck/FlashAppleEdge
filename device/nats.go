package device

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func (se *State) initNataSrv() (err error) {

	// Create a new NATS server configuration
	opts := &server.Options{
		Host:     "0.0.0.0",
		Port:     Nats_port,
		HTTPPort: Nats_port + 1,
	}

	// Create a new NATS server instance
	se.NatsServer, err = server.NewServer(opts)
	if err != nil {
		return fmt.Errorf("Error creating NATS server: %v", err)
	}
	fmt.Printf("NATS server created at %d", Nats_port)

	se.NatsServer.Start()

	if !se.NatsServer.ReadyForConnections(4 * time.Second) {
		return fmt.Errorf("Nats Server konnte nicht nach 4 sek starten")
	}
	fmt.Printf("NATS server is ready for connections at %s:%d", opts.Host, Nats_port)

	if se.NatsClient, err = nats.Connect(se.NatsServer.ClientURL()); err != nil {
		return fmt.Errorf("Error connecting NATS client: %v", err)
	}

	return nil
}

func (se *State) chanSubWrapper(subj string, nmsg chan *nats.Msg) (*nats.Subscription, error) {

	return se.NatsClient.ChanSubscribe(subj, nmsg)
}

func (se *State) chanLatestMsg(ch chan *nats.Msg) (*nats.Subscription, error) {
	topic := se.Cfg.Base + ".>"
	fmt.Println("ChanLatestMsg called()")
	return se.NatsClient.ChanSubscribe(topic, ch)
}
