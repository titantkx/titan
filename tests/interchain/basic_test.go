package interchain_test

import (
	"context"
	"testing"

	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/conformance"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"go.uber.org/zap/zaptest"
)

func TestConformance(t *testing.T) {
	var (
		ctx                        = context.Background()
		rep                        = testreporter.NewNopReporter()
		chainIDTitan, chainIDGaia1 = "titan_18887-1", "cosmoshub-1"
	)

	numValidators := 1 // Defines how many validators should be used in each network.
	numFullNodes := 0  // Defines how many additional full nodes should be used in each network.

	// Here we define our ChainFactory by instantiating a new instance of the BuiltinChainFactory exposed in interchaintest.
	// We use the ChainSpec type to fully describe which chains we want to use in our tests.
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		titanChainSpec(ctx, chainIDTitan, numValidators, numFullNodes),
		{
			Name:          "gaia",
			ChainName:     chainIDGaia1,
			Version:       "v13.0.1",
			NumValidators: &numValidators,
			NumFullNodes:  &numFullNodes,
		},
	})

	rlyFactory := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	)

	// Test will now run the conformance test suite against both of our chains, ensuring that they both have basic
	// IBC capabilities properly implemented and work with both the Go relayer and Hermes.
	conformance.Test(t, ctx, []interchaintest.ChainFactory{cf}, []interchaintest.RelayerFactory{rlyFactory}, rep)
}
