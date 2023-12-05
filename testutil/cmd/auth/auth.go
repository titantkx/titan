package auth

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
)

type AccountsResponse struct {
	Accounts   []Account  `json:"accounts"`
	Pagination Pagination `json:"pagination"`
}

type Account struct {
	Type          string              `json:"@type"`
	Address       string              `json:"address"`
	PubKey        *testutil.PublicKey `json:"pub_key"`
	AccountNumber testutil.Int        `json:"account_number"`
	Sequence      testutil.Int        `json:"sequence"`
	BaseAccount   *BaseAccount        `json:"base_account"`
	CodeHash      string              `json:"code_hash"`
}

func (a *Account) GetAddress() string {
	if a.BaseAccount != nil {
		return a.BaseAccount.Address
	}
	return a.Address
}

func (a Account) GetPubKey() *testutil.PublicKey {
	if a.BaseAccount != nil {
		return a.BaseAccount.PubKey
	}
	return a.PubKey
}

func (a Account) GetAccountNumber() int64 {
	if a.BaseAccount != nil {
		return a.BaseAccount.AccountNumber.Int64()
	}
	return a.AccountNumber.Int64()
}

func (a Account) GetSequence() int64 {
	if a.BaseAccount != nil {
		return a.BaseAccount.Sequence.Int64()
	}
	return a.Sequence.Int64()
}

type BaseAccount struct {
	Address       string              `json:"address"`
	PubKey        *testutil.PublicKey `json:"pub_key"`
	AccountNumber testutil.Int        `json:"account_number"`
	Sequence      testutil.Int        `json:"sequence"`
}

type Pagination struct {
	NextKey string       `json:"next_key"`
	Total   testutil.Int `json:"total"`
}

func MustGetAccounts(t testing.TB, height int64) <-chan Account {
	ch := make(chan Account, 100)
	go func() {
		defer close(ch)
		limit := 100
		offset := 0
		for {
			args := []string{
				"auth",
				"accounts",
				"--limit=" + strconv.Itoa(limit),
				"--offset=" + strconv.Itoa(offset),
			}
			if height > 0 {
				args = append(args, "--height="+testutil.FormatInt(height))
			}

			var resp AccountsResponse
			cmd.MustQuery(t, &resp, args...)

			for _, account := range resp.Accounts {
				require.NotEmpty(t, account.Type)
				require.NotEmpty(t, account.GetAddress())
				if pk := account.GetPubKey(); pk != nil {
					require.NotEmpty(t, pk.Type)
					require.NotEmpty(t, pk.Key)
				}
				require.GreaterOrEqual(t, account.GetAccountNumber(), int64(0))
				require.GreaterOrEqual(t, account.GetSequence(), int64(0))
				ch <- account
			}

			if resp.Pagination.NextKey == "" {
				break
			}
			offset += limit
		}
	}()
	return ch
}
