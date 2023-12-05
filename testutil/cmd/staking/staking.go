package staking

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
)

const (
	BOND_STATUS_BONDED    = "BOND_STATUS_BONDED"
	BOND_STATUS_UNBONDED  = "BOND_STATUS_UNBONDED"
	BOND_STATUS_UNBONDING = "BOND_STATUS_UNBONDING"
)

type Validator struct {
	OperatorAddress   string             `json:"operator_address"`
	ConsensusPubkey   testutil.PublicKey `json:"consensus_pubkey"`
	Commission        Commission         `json:"commission"`
	MinSelfDelegation testutil.Int       `json:"min_self_delegation"`
	Jailed            bool               `json:"jailed"`
	Status            string             `json:"status"`
	Tokens            testutil.BigInt    `json:"tokens"`
	DelegatorShares   testutil.BigFloat  `json:"delegator_shares"`
}

type Commission struct {
	CommissionRates CommissionRates `json:"commission_rates"`
}

type CommissionRates struct {
	Rate          testutil.Float `json:"rate"`
	MaxRate       testutil.Float `json:"max_rate"`
	MaxChangeRate testutil.Float `json:"max_change_rate"`
}

func MustGetValidator(t testing.TB, address string) Validator {
	var val Validator
	cmd.MustQuery(t, &val, "staking", "validator", address)
	require.Equal(t, address, val.OperatorAddress)
	return val
}

type DelegationResponse struct {
	Delegation Delegation `json:"delegation"`
	Balance    Balance    `json:"balance"`
}

type Delegation struct {
	DelegatorAddress string            `json:"delegator_address"`
	ValidatorAddress string            `json:"validator_address"`
	Shares           testutil.BigFloat `json:"shares"`
}

type Balance struct {
	Denom  string          `json:"denom"`
	Amount testutil.BigInt `json:"amount"`
}

func MustGetDelegation(t testing.TB, delegator string, validator string) DelegationResponse {
	var resp DelegationResponse
	cmd.MustQuery(t, &resp, "staking", "delegation", delegator, validator)
	require.Equal(t, delegator, resp.Delegation.DelegatorAddress)
	require.Equal(t, validator, resp.Delegation.ValidatorAddress)
	return resp
}

func MustCreateValidator(t testing.TB, valPk testutil.PublicKey, amount string, commissionRate float64, commissionMaxRate float64, commissionMaxChangeRate float64, minSelfDelegation int64, from string) Validator {
	balBefore := bank.MustGetBalance(t, from, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "create-validator", "--pubkey="+valPk.String(), "--amount="+amount, "--commission-rate="+testutil.FormatFloat(commissionRate), "--commission-max-rate="+testutil.FormatFloat(commissionMaxRate), "--commission-max-change-rate="+testutil.FormatFloat(commissionMaxChangeRate), "--min-self-delegation="+testutil.FormatInt(minSelfDelegation), "--from="+from)

	balAfter := bank.MustGetBalance(t, from, "utkx", 0)

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	stakedAmount := testutil.MustGetUtkxAmount(t, amount)

	require.Equal(t, balBefore.Sub(coinSpent).Sub(stakedAmount), balAfter)

	var valAddr string
	var actualStakedAmount testutil.BigInt

	for _, event := range tx.Events {
		if event.Type == "create_validator" {
			for _, att := range event.Attributes {
				if att.Key == "validator" {
					valAddr = att.Value
				} else if att.Key == "amount" {
					actualStakedAmount = testutil.MustGetUtkxAmount(t, att.Value)
				}
			}
		}
	}

	require.NotEmpty(t, valAddr)
	require.False(t, actualStakedAmount.IsZero())
	require.Equal(t, stakedAmount, actualStakedAmount)

	val := MustGetValidator(t, valAddr)

	require.Equal(t, valPk.Type, val.ConsensusPubkey.Type)
	require.Equal(t, valPk.Key, val.ConsensusPubkey.Key)
	require.Equal(t, commissionRate, val.Commission.CommissionRates.Rate.Float64())
	require.Equal(t, commissionMaxRate, val.Commission.CommissionRates.MaxRate.Float64())
	require.Equal(t, commissionMaxChangeRate, val.Commission.CommissionRates.MaxChangeRate.Float64())
	require.Equal(t, minSelfDelegation, val.MinSelfDelegation.Int64())
	require.False(t, val.Jailed)
	require.Equal(t, val.Status, BOND_STATUS_BONDED)
	require.Equal(t, stakedAmount, val.Tokens)

	del := MustGetDelegation(t, from, valAddr)

	require.Equal(t, stakedAmount, del.Balance.Amount)

	return val
}

func MustDelegate(t testing.TB, valAddr string, amount string, from string) {
	valBefore := MustGetValidator(t, valAddr)
	balBefore := bank.MustGetBalance(t, from, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "delegate", valAddr, amount, "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	balAfter := bank.MustGetBalance(t, from, "utkx", 0)

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	delegatedAmount := testutil.MustGetUtkxAmount(t, amount)

	require.Equal(t, balBefore.Sub(coinSpent).Sub(delegatedAmount), balAfter)
	require.Equal(t, valBefore.Tokens.Add(delegatedAmount), valAfter.Tokens)

	del := MustGetDelegation(t, from, valAddr)

	require.Equal(t, delegatedAmount, del.Balance.Amount)
}

func MustRedelegate(t testing.TB, srcVal string, dstVal, amount string, from string) {
	srcValBefore := MustGetValidator(t, srcVal)
	dstValBefore := MustGetValidator(t, dstVal)
	balBefore := bank.MustGetBalance(t, from, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "redelegate", srcVal, dstVal, amount, "--from="+from)

	srcValAfter := MustGetValidator(t, srcVal)
	dstValAfter := MustGetValidator(t, dstVal)
	balAfter := bank.MustGetBalance(t, from, "utkx", 0)

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	redelegatedAmount := testutil.MustGetUtkxAmount(t, amount)
	reward := mustGetReward(t, tx.Events)

	require.Equal(t, balBefore.Sub(coinSpent).Add(reward), balAfter)
	require.Equal(t, srcValBefore.Tokens.Sub(redelegatedAmount), srcValAfter.Tokens)
	require.Equal(t, dstValBefore.Tokens.Add(redelegatedAmount), dstValAfter.Tokens)

	del := MustGetDelegation(t, from, dstVal)

	require.Equal(t, redelegatedAmount, del.Balance.Amount)
}

func MustUnbond(t testing.TB, valAddr string, amount string, from string) txcmd.Tx {
	valBefore := MustGetValidator(t, valAddr)
	balBefore := bank.MustGetBalance(t, from, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "unbond", valAddr, amount, "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	balAfter := bank.MustGetBalance(t, from, "utkx", 0)

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	unbondedAmount := testutil.MustGetUtkxAmount(t, amount)
	reward := mustGetReward(t, tx.Events)

	require.Equal(t, balBefore.Sub(coinSpent).Add(reward), balAfter)
	require.Equal(t, valBefore.Tokens.Sub(unbondedAmount), valAfter.Tokens)

	return tx
}

func MustCancelUnbound(t testing.TB, valAddr string, amount string, creationHeight int64, from string) {
	valBefore := MustGetValidator(t, valAddr)
	balBefore := bank.MustGetBalance(t, from, "utkx", 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "cancel-unbond", valAddr, amount, testutil.FormatInt(creationHeight), "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	balAfter := bank.MustGetBalance(t, from, "utkx", 0)

	coinSpent := tx.GasWanted.Mul(testutil.MakeBigInt(10)) // Gas price == 10 utkx
	unbondedAmount := testutil.MustGetUtkxAmount(t, amount)
	reward := mustGetReward(t, tx.Events)

	require.Equal(t, balBefore.Sub(coinSpent).Add(reward), balAfter)
	require.Equal(t, valBefore.Tokens.Add(unbondedAmount), valAfter.Tokens)
}

func mustGetReward(t testing.TB, events []txcmd.Event) testutil.BigInt {
	reward := testutil.MakeBigInt(0)
	for _, event := range events {
		if event.Type == "withdraw_rewards" {
			for _, att := range event.Attributes {
				if att.Key == "amount" {
					reward = testutil.MustGetUtkxAmount(t, att.Value)
				}
			}
		}
	}
	return reward
}
