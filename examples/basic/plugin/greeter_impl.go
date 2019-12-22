package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	example "github.com/hashicorp/go-plugin/examples/basic/commons"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Here is a real implementation of Greeter
type GreeterHello struct {
	logger hclog.Logger
}

func (g *GreeterHello) Greet(c *example.Client) string {
	g.logger.Debug("message from GreeterHello.Greet")

	attribute := c.WalkAttribute()
	var val string
	if err := c.EvaluateExpr(attribute.Expr, &val); err != nil {
		panic(err)
	}

	return fmt.Sprintf("Hello! %s", val)
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	gob.Register(&hclsyntax.TemplateExpr{})
	gob.Register(&hclsyntax.LiteralValueExpr{})
	gob.Register(&hclsyntax.ScopeTraversalExpr{})
	gob.Register(hcl.TraverseRoot{})
	gob.Register(hcl.TraverseAttr{})

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	greeter := &GreeterHello{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &example.GreeterPlugin{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
