package example

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
type Greeter interface {
	Greet(*Client) string
}

// Here is an implementation that talks over RPC
type GreeterRPC struct {
	client *rpc.Client
	broker *plugin.MuxBroker
}

func (g *GreeterRPC) Greet(c *Client) string {
	brokerID := g.broker.NextId()
	go g.broker.AcceptAndServe(brokerID, &Server{})

	var resp string
	err := g.client.Call("Plugin.Greet", brokerID, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type GreeterRPCServer struct {
	// This is the real implementation
	Impl   Greeter
	Broker *plugin.MuxBroker
}

func (s *GreeterRPCServer) Greet(brokerID uint32, resp *string) error {
	conn, err := s.Broker.Dial(brokerID)
	if err != nil {
		panic(err)
	}
	*resp = s.Impl.Greet(&Client{rpcClient: rpc.NewClient(conn)})
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type GreeterPlugin struct {
	// Impl Injection
	Impl Greeter
}

func (p *GreeterPlugin) Server(b *plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl, Broker: b}, nil
}

func (GreeterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GreeterRPC{client: c, broker: b}, nil
}
