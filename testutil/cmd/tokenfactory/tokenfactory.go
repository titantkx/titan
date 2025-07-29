package tokenfactory

import (
	"context"
	"fmt"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/titantkx/titan/testutil"
	txcmd "github.com/titantkx/titan/testutil/cmd/tx"
)

func MustCreateNewToken(t testutil.TestingT, creator string, denom string) string {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{
		"tokenfactory",
		"create-denom",
		denom,
		"--from=" + creator,
	}

	tx := txcmd.MustExecTx(t, ctx, args...)

	creatorInEvent := strings.Trim(tx.MustGetEventAttributeValue(t, "create_denom", "creator"), "\"")
	tokenDenom := strings.Trim(tx.MustGetEventAttributeValue(t, "create_denom", "new_token_denom"), "\"")
	require.Equal(t, creator, creatorInEvent)
	require.Equal(t, fmt.Sprintf("factory/%s/%s", creator, denom), tokenDenom)

	return tokenDenom
}

func MustMintTokenFactory(t testutil.TestingT, creator string, receiver string, tokenDenom string, amount testutil.Int) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	args := []string{
		"tokenfactory",
		"mint",
		fmt.Sprintf("%s%s", amount.String(), tokenDenom),
		receiver,
		"--from=" + creator,
	}

	tx := txcmd.MustExecTx(t, ctx, args...)

	mintToAddr := strings.Trim(tx.MustGetEventAttributeValue(t, "tf_mint", "mint_to_address"), "\"")
	amountDenom := strings.Trim(tx.MustGetEventAttributeValue(t, "tf_mint", "amount"), "\"")
	require.Equal(t, receiver, mintToAddr)
	require.Equal(t, fmt.Sprintf("%s%s", amount.String(), tokenDenom), amountDenom)
}
