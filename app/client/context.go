package client

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client"
)

type Context struct {
	client.Context
	Deadline time.Time
}

func (ctx Context) WithDeadline(deadline time.Time) Context {
	ctx.Deadline = deadline
	return ctx
}
