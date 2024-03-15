package client

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servercmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewTxFlags() *pflag.FlagSet {
	var cmd cobra.Command
	flags.AddTxFlagsToCmd(&cmd)
	return cmd.Flags()
}

func NewQueryFlags() *pflag.FlagSet {
	var cmd cobra.Command
	flags.AddQueryFlagsToCmd(&cmd)
	return cmd.Flags()
}

func NewPaginationFlags(query string) *pflag.FlagSet {
	var cmd cobra.Command
	flags.AddPaginationFlagsToCmd(&cmd, query)
	flags.AddQueryFlagsToCmd(&cmd)
	return cmd.Flags()
}

func ReadTxFlags(ctx Context, flags *pflag.FlagSet) (Context, error) {
	var cmd cobra.Command
	cmd.Flags().AddFlagSet(flags)
	cmd.SetContext(servercmd.CreateExecuteContext(context.Background()))
	if err := client.SetCmdClientContext(&cmd, ctx.Context); err != nil {
		return ctx, err
	}
	clientCtx, err := client.GetClientTxContext(&cmd)
	if err != nil {
		return ctx, err
	}
	ctx.Context = clientCtx
	return ctx, nil
}
