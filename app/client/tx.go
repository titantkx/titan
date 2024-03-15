package client

import (
	"errors"
	"math/rand"
	"time"

	"github.com/spf13/pflag"

	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
)

func GenerateTx(ctx Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) ([]byte, error) {
	txf, err := clienttx.NewFactoryCLI(ctx.Context, flagSet)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	return GenerateTxWithFactory(ctx, txf, msgs...)
}

func GenerateTxWithFactory(ctx Context, txf clienttx.Factory, msgs ...sdk.Msg) ([]byte, error) {
	if txf.SimulateAndExecute() {
		if ctx.Offline {
			return nil, errors.New("cannot estimate gas in offline mode")
		}

		// Prepare TxFactory with acc & seq numbers as CalculateGas requires
		// account and sequence numbers to be set
		preparedTxf, err := txf.Prepare(ctx.Context)
		if err != nil {
			return nil, err
		}

		_, adjusted, err := clienttx.CalculateGas(ctx, preparedTxf, msgs...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(adjusted)
	}

	unsignedTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	return ctx.TxConfig.TxJSONEncoder()(unsignedTx.GetTx())
}

func BroadcastTx(ctx Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf, err := clienttx.NewFactoryCLI(ctx.Context, flagSet)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	return BroadcastTxWithFactory(ctx, txf, msgs...)
}

func BroadcastTxWithFactory(ctx Context, txf clienttx.Factory, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf, err := txf.Prepare(ctx.Context)
	if err != nil {
		return nil, err
	}

	if txf.SimulateAndExecute() || ctx.Simulate {
		if ctx.Offline {
			return nil, errors.New("cannot estimate gas in offline mode")
		}

		_, adjusted, err := clienttx.CalculateGas(ctx, txf, msgs...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(adjusted)
	}

	if ctx.Simulate {
		return nil, nil
	}

	tx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	err = clienttx.Sign(txf, ctx.GetFromName(), tx, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := ctx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return nil, err
	}

	// broadcast to a Tendermint node
	return ctx.BroadcastTx(txBytes)
}

func QueryTx(ctx Context, txHash string) (*sdk.TxResponse, error) {
	timeout := time.NewTimer(time.Until(ctx.Deadline))
	defer timeout.Stop()

	for retry := 1; ; retry++ {
		tx, err := authtx.QueryTx(ctx.Context, txHash)
		if err != nil {
			// If transaction not found, wait until it is delivered or timeout
			if IsNotFound(err) {
				// Exponential retry backoff with 10% jitter
				backoff := time.Duration(float64(2*retry) * float64(time.Second) * (0.9 + rand.Float64()*0.2))
				select {
				case <-time.After(backoff):
					continue
				case <-timeout.C:
					return nil, ErrDeadlineExceeded
				}
			}
			return nil, err
		}
		return tx, nil
	}
}

func BroadcastAndQueryTx(ctx Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	tx, err := BroadcastTx(ctx, flagSet, msgs...)
	if err != nil {
		return nil, err
	}

	return QueryTx(ctx, tx.TxHash)
}
