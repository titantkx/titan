package client

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	etherminthd "github.com/titantkx/ethermint/crypto/hd"
	etherminttypes "github.com/titantkx/ethermint/types"
	"github.com/titantkx/titan/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateClientContext(conf Config) (Context, error) {
	encodingConfig := app.MakeEncodingConfig()
	ctx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithChainID(conf.GetChainID()).
		WithHomeDir(conf.GetHomeDir()).
		WithKeyringDir(conf.GetKeyringDir()).
		WithOutputFormat(conf.GetOutputFormat()).
		// Ethermint
		WithKeyringOptions(etherminthd.EthSecp256k1Option())

	keyring, err := client.NewKeyringFromBackend(ctx, conf.GetKeyringBackend())
	if err != nil {
		return Context{}, fmt.Errorf("couldn't get key ring: %v", err)
	}

	ctx = ctx.WithKeyring(keyring)

	client, err := client.NewClientFromNode(conf.GetNodeURI())
	if err != nil {
		return Context{}, fmt.Errorf("couldn't get client from nodeURI: %v", err)
	}

	ctx = ctx.WithNodeURI(conf.GetNodeURI()).
		WithClient(client).
		WithBroadcastMode(conf.GetBroadcastMode())

	if conf.GRPC != "" {
		var dialOpts []grpc.DialOption

		if conf.GRPCInsecure {
			dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		} else {
			dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				MinVersion: tls.VersionTLS12,
			})))
		}

		grpcClient, err := grpc.Dial(conf.GRPC, dialOpts...)
		if err != nil {
			return Context{}, err
		}
		ctx = ctx.WithGRPCClient(grpcClient)
	}

	ctx, err = addDefaultKeyInfoKeyring(ctx)
	if err != nil {
		return Context{}, err
	}

	return Context{Context: ctx}, nil
}

// if in `--dry-run` mode and clientCtx.Keyring is `keyring.BackendMemory` and empty, add default key info
// it will prevent error `cannot build signature for simulation, key records slice is empty` while running `--dry-run`
func addDefaultKeyInfoKeyring(clientCtx client.Context) (client.Context, error) {
	if clientCtx.Simulate && clientCtx.Keyring != nil && clientCtx.Keyring.Backend() == keyring.BackendMemory {
		kr := clientCtx.Keyring
		records, _ := kr.List()
		if len(records) == 0 {
			// add default key info
			_, _, err := kr.NewMnemonic("foo", keyring.English, etherminttypes.BIP44HDPath, keyring.DefaultBIP39Passphrase, etherminthd.EthSecp256k1)
			if err != nil {
				return clientCtx, err
			}
		}

		records, _ = kr.List()

		if len(records) == 0 {
			return clientCtx, errors.New("cannot add default key info")
		}
	}

	return clientCtx, nil
}
