package example

import (
	"log"
	"net/rpc"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Client struct {
	rpcClient *rpc.Client
}

func (c *Client) WalkAttribute() hcl.Attribute {
	var resp hcl.Attribute
	if err := c.rpcClient.Call("Plugin.Walk", new(interface{}), &resp); err != nil {
		log.Print(err)
	}

	return resp
}

func (c *Client) EvaluateExpr(expr hcl.Expression, ret interface{}) error {
	var resp cty.Value
	if err := c.rpcClient.Call("Plugin.EvaluateExpr", &expr, &resp); err != nil {
		log.Print(err)
	}

	if err := gocty.FromCtyValue(resp, ret); err != nil {
		log.Print(err)
	}

	return nil
}
