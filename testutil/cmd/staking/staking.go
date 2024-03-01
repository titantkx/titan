package staking

import (
	"context"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
	"github.com/tokenize-titan/titan/testutil/cmd"
	"github.com/tokenize-titan/titan/testutil/cmd/bank"
	"github.com/tokenize-titan/titan/testutil/cmd/slashing"
	txcmd "github.com/tokenize-titan/titan/testutil/cmd/tx"
	"github.com/tokenize-titan/titan/utils"
)

const (
	BOND_STATUS_BONDED    = "BOND_STATUS_BONDED"
	BOND_STATUS_UNBONDED  = "BOND_STATUS_UNBONDED"
	BOND_STATUS_UNBONDING = "BOND_STATUS_UNBONDING"
)

type Validator struct {
	OperatorAddress   string                   `json:"operator_address"`
	ConsensusPubkey   testutil.SinglePublicKey `json:"consensus_pubkey"`
	Commission        Commission               `json:"commission"`
	MinSelfDelegation testutil.Int             `json:"min_self_delegation"`
	Jailed            bool                     `json:"jailed"`
	Status            string                   `json:"status"`
	Tokens            testutil.Int             `json:"tokens"`
	DelegatorShares   testutil.Float           `json:"delegator_shares"`
}

type Commission struct {
	CommissionRates CommissionRates `json:"commission_rates"`
}

type CommissionRates struct {
	Rate          testutil.Float `json:"rate"`
	MaxRate       testutil.Float `json:"max_rate"`
	MaxChangeRate testutil.Float `json:"max_change_rate"`
}

type StakingParams struct {
	BondDenom               string            `json:"bond_denom"`
	HistoricalEntries       int64             `json:"historical_entries"`
	MaxEntries              int64             `json:"max_entries"`
	MaxValidators           int64             `json:"max_validators"`
	MinCommissionRate       testutil.Float    `json:"min_commission_rate"`
	UnbondingTime           testutil.Duration `json:"unbonding_time"`
	GlobalMinSelfDelegation testutil.Int      `json:"global_min_self_delegation"`
}

func MustGetValidator(t testutil.TestingT, address string) Validator {
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
	DelegatorAddress string         `json:"delegator_address"`
	ValidatorAddress string         `json:"validator_address"`
	Shares           testutil.Float `json:"shares"`
}

type Balance struct {
	Denom  string       `json:"denom"`
	Amount testutil.Int `json:"amount"`
}

func GetDelegation(delegator string, validator string) (*DelegationResponse, error) {
	var resp DelegationResponse
	err := cmd.Query(&resp, "staking", "delegation", delegator, validator)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func MustGetDelegation(t testutil.TestingT, delegator string, validator string) DelegationResponse {
	var resp DelegationResponse
	cmd.MustQuery(t, &resp, "staking", "delegation", delegator, validator)
	require.Equal(t, delegator, resp.Delegation.DelegatorAddress)
	require.Equal(t, validator, resp.Delegation.ValidatorAddress)
	return resp
}

func MustCreateValidator(t testutil.TestingT, valPk testutil.SinglePublicKey, amount string, commissionRate float64, commissionMaxRate float64, commissionMaxChangeRate float64, minSelfDelegation testutil.Int, from string) Validator {
	balBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "create-validator", "--pubkey="+valPk.String(), "--amount="+amount, "--commission-rate="+testutil.FormatFloat(commissionRate), "--commission-max-rate="+testutil.FormatFloat(commissionMaxRate), "--commission-max-change-rate="+testutil.FormatFloat(commissionMaxChangeRate), "--min-self-delegation="+minSelfDelegation.String(), "--from="+from)

	balAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	coinSpent := tx.MustGetDeductFeeAmount(t)
	stakedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	sharedAmount := stakedAmount.Float()

	require.Equal(t, balBefore.Sub(coinSpent).Sub(stakedAmount), balAfter)

	valAddr := tx.MustGetEventAttributeValue(t, "create_validator", "validator")
	actualStakedAmount := mustGetStakedAmount(t, tx)

	require.NotEmpty(t, valAddr)
	require.False(t, actualStakedAmount.IsZero())
	require.Equal(t, stakedAmount, actualStakedAmount)

	val := MustGetValidator(t, valAddr)

	require.Equal(t, valPk.Type, val.ConsensusPubkey.Type)
	require.Equal(t, valPk.Key, val.ConsensusPubkey.Key)
	require.Equal(t, commissionRate, val.Commission.CommissionRates.Rate.Float64())
	require.Equal(t, commissionMaxRate, val.Commission.CommissionRates.MaxRate.Float64())
	require.Equal(t, commissionMaxChangeRate, val.Commission.CommissionRates.MaxChangeRate.Float64())
	require.Equal(t, minSelfDelegation, val.MinSelfDelegation)
	require.False(t, val.Jailed)
	require.Equal(t, stakedAmount, val.Tokens)
	require.Equal(t, sharedAmount.String(), val.DelegatorShares.String())

	del := MustGetDelegation(t, from, valAddr)

	require.Equal(t, stakedAmount, del.Balance.Amount)
	require.Equal(t, sharedAmount.String(), del.Delegation.Shares.String())

	return val
}

func MustErrCreateValidator(t testutil.TestingT, expErr string, valPk testutil.SinglePublicKey, amount string, commissionRate float64, commissionMaxRate float64, commissionMaxChangeRate float64, minSelfDelegation testutil.Int, from string) {
	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	txcmd.MustErrExecTx(t, ctx, expErr, "staking", "create-validator", "--pubkey="+valPk.String(), "--amount="+amount, "--commission-rate="+testutil.FormatFloat(commissionRate), "--commission-max-rate="+testutil.FormatFloat(commissionMaxRate), "--commission-max-change-rate="+testutil.FormatFloat(commissionMaxChangeRate), "--min-self-delegation="+minSelfDelegation.String(), "--from="+from)
}

func MustDelegate(t testutil.TestingT, valAddr string, amount string, from string) {
	valBefore := MustGetValidator(t, valAddr)
	balBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "delegate", valAddr, amount, "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	balAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	coinSpent := tx.MustGetDeductFeeAmount(t)
	delegatedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	slashedAmount := mustGetSlashedAmount(t, valBefore, valAfter)
	sharedAmount := valBefore.DelegatorShares.Mul(delegatedAmount.DivFloat(valBefore.Tokens))

	require.Equal(t, balBefore.Sub(coinSpent).Sub(delegatedAmount), balAfter)
	require.Equal(t, valBefore.Tokens.Sub(slashedAmount).Add(delegatedAmount), valAfter.Tokens)
	require.Equal(t, valBefore.DelegatorShares.Add(sharedAmount).String(), valAfter.DelegatorShares.String())

	del := MustGetDelegation(t, from, valAddr)

	require.Equal(t, delegatedAmount, del.Balance.Amount)
	require.Equal(t, sharedAmount.String(), del.Delegation.Shares.String())
}

func MustRedelegate(t testutil.TestingT, srcVal string, dstVal, amount string, from string) {
	srcValBefore := MustGetValidator(t, srcVal)
	dstValBefore := MustGetValidator(t, dstVal)
	balBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)
	srcDelBefore := MustGetDelegation(t, from, srcVal)
	dstDelBefore, err := GetDelegation(from, dstVal)

	if err != nil {
		require.ErrorContains(t, err, "NotFound")
	} else {
		require.NotNil(t, dstDelBefore)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "redelegate", srcVal, dstVal, amount, "--from="+from)

	srcValAfter := MustGetValidator(t, srcVal)
	dstValAfter := MustGetValidator(t, dstVal)
	balAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	reward := mustGetReward(t, tx)
	coinSpent := tx.MustGetDeductFeeAmount(t)
	redelegatedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	srcSlashedAmount := mustGetSlashedAmount(t, srcValBefore, srcValAfter)
	dstSlashedAmount := mustGetSlashedAmount(t, dstValBefore, dstValAfter)
	unbondedShares := srcValBefore.DelegatorShares.Mul(redelegatedAmount.DivFloat(srcValBefore.Tokens))
	newShares := dstValBefore.DelegatorShares.Mul(redelegatedAmount.DivFloat(dstValBefore.Tokens))

	require.Equal(t, balBefore.Sub(coinSpent).Add(reward), balAfter)
	require.Equal(t, srcValBefore.Tokens.Sub(srcSlashedAmount).Sub(redelegatedAmount), srcValAfter.Tokens)
	require.Equal(t, dstValBefore.Tokens.Sub(dstSlashedAmount).Add(redelegatedAmount), dstValAfter.Tokens)
	require.Equal(t, srcValBefore.DelegatorShares.Sub(unbondedShares).String(), srcValAfter.DelegatorShares.String())
	require.Equal(t, dstValBefore.DelegatorShares.Add(newShares).String(), dstValAfter.DelegatorShares.String())

	srcDelAfter, err := GetDelegation(from, srcVal)

	if redelegatedAmount.Cmp(srcDelBefore.Balance.Amount) == 0 {
		require.Error(t, err)
		require.ErrorContains(t, err, "NotFound")
	} else {
		require.NoError(t, err)
		require.NotNil(t, srcDelAfter)
		require.Equal(t, srcDelBefore.Balance.Amount.Sub(redelegatedAmount), srcDelAfter.Balance.Amount)
		require.Equal(t, srcDelBefore.Delegation.Shares.Sub(unbondedShares).String(), srcDelAfter.Delegation.Shares.String())
	}

	dstDelAfter := MustGetDelegation(t, from, dstVal)

	if dstDelBefore == nil {
		require.Equal(t, redelegatedAmount, dstDelAfter.Balance.Amount)
		require.Equal(t, newShares.String(), dstDelAfter.Delegation.Shares.String())
	} else {
		require.Equal(t, dstDelBefore.Balance.Amount.Add(redelegatedAmount), dstDelAfter.Balance.Amount)
		require.Equal(t, dstDelBefore.Delegation.Shares.Add(newShares).String(), dstDelAfter.Delegation.Shares.String())
	}
}

func MustUnbond(t testutil.TestingT, valAddr string, amount string, from string) txcmd.TxResponse {
	valBefore := MustGetValidator(t, valAddr)
	balBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)
	delBefore := MustGetDelegation(t, from, valAddr)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "unbond", valAddr, amount, "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	balAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	reward := mustGetReward(t, tx)
	coinSpent := tx.MustGetDeductFeeAmount(t)
	unbondedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	slashedAmount := mustGetSlashedAmount(t, valBefore, valAfter)
	unbondedShares := valBefore.DelegatorShares.Mul(unbondedAmount.DivFloat(valBefore.Tokens))

	require.Equal(t, balBefore.Sub(coinSpent).Add(reward), balAfter)
	require.Equal(t, valBefore.Tokens.Sub(slashedAmount).Sub(unbondedAmount), valAfter.Tokens)
	require.Equal(t, valBefore.DelegatorShares.Sub(unbondedShares).String(), valAfter.DelegatorShares.String())

	delAfter, err := GetDelegation(from, valAddr)

	if unbondedAmount.Cmp(delBefore.Balance.Amount) == 0 {
		require.Error(t, err)
		require.ErrorContains(t, err, "NotFound")
	} else {
		require.NoError(t, err)
		require.NotNil(t, delAfter)
		require.Equal(t, delBefore.Balance.Amount.Sub(unbondedAmount), delAfter.Balance.Amount)
		require.Equal(t, delBefore.Delegation.Shares.Sub(unbondedShares).String(), delAfter.Delegation.Shares.String())
	}

	return tx
}

func MustCancelUnbound(t testutil.TestingT, valAddr string, amount string, creationHeight int64, from string) {
	valBefore := MustGetValidator(t, valAddr)
	balBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)
	delBefore, err := GetDelegation(from, valAddr)

	if err != nil {
		require.ErrorContains(t, err, "NotFound")
	} else {
		require.NotNil(t, delBefore)
	}

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "cancel-unbond", valAddr, amount, testutil.FormatInt(creationHeight), "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	balAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	reward := mustGetReward(t, tx)
	coinSpent := tx.MustGetDeductFeeAmount(t)
	unbondedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	slashedAmount := mustGetSlashedAmount(t, valBefore, valAfter)
	unbondedShares := valBefore.DelegatorShares.Mul(unbondedAmount.DivFloat(valBefore.Tokens))

	require.Equal(t, balBefore.Sub(coinSpent).Add(reward), balAfter)
	require.Equal(t, valBefore.Tokens.Sub(slashedAmount).Add(unbondedAmount), valAfter.Tokens)
	require.Equal(t, valBefore.DelegatorShares.Add(unbondedShares).String(), valAfter.DelegatorShares.String())

	delAfter := MustGetDelegation(t, from, valAddr)

	if delBefore == nil {
		require.Equal(t, unbondedAmount, delAfter.Balance.Amount)
		require.Equal(t, unbondedShares.String(), delAfter.Delegation.Shares.String())
	} else {
		require.Equal(t, delBefore.Balance.Amount.Add(unbondedAmount), delAfter.Balance.Amount)
		require.Equal(t, delBefore.Delegation.Shares.Add(unbondedShares).String(), delAfter.Delegation.Shares.String())
	}
}

func mustGetReward(t testutil.TestingT, txr txcmd.TxResponse) testutil.Int {
	reward := testutil.MakeInt(0)
	for _, event := range txr.Events {
		if event.Type == "withdraw_rewards" {
			attr := event.FindAttribute("amount")
			require.NotNil(t, attr)
			reward = reward.Add(testutil.MustGetBaseDenomAmount(t, attr.Value))
		}
	}
	return reward
}

func mustGetStakedAmount(t testutil.TestingT, txr txcmd.TxResponse) testutil.Int {
	amount := txr.MustGetEventAttributeValue(t, "create_validator", "amount")
	return testutil.MustGetBaseDenomAmount(t, amount)
}

func mustGetSlashedAmount(t testutil.TestingT, valBefore Validator, valAfter Validator) testutil.Int {
	if !valAfter.Jailed || valBefore.Jailed {
		// Validator is not jailed or was jailed before
		return testutil.MakeInt(0)
	}
	if valAfter.Tokens.Cmp(valAfter.MinSelfDelegation) < 0 {
		// Validator was jailed due to self delegation lower than min self delegation
		return testutil.MakeInt(0)
	}
	params := slashing.MustGetParams(t)
	return valBefore.Tokens.Float().Mul(params.SlashFractionDowntime).Int()
}

func MustGetParams(t testutil.TestingT) StakingParams {
	var params StakingParams
	cmd.MustQuery(t, &params, "staking", "params")
	return params
}

func MustCreateValidatorForOther(t testutil.TestingT, valPk testutil.SinglePublicKey, amount string, commissionRate float64, commissionMaxRate float64, commissionMaxChangeRate float64, minSelfDelegation testutil.Int, from string, delAddr string) Validator {
	delBalBefore := bank.MustGetBalance(t, delAddr, utils.BaseDenom, 0)
	fromBalBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "create-validator-for-other", delAddr, "--pubkey="+valPk.String(), "--amount="+amount, "--commission-rate="+testutil.FormatFloat(commissionRate), "--commission-max-rate="+testutil.FormatFloat(commissionMaxRate), "--commission-max-change-rate="+testutil.FormatFloat(commissionMaxChangeRate), "--min-self-delegation="+minSelfDelegation.String(), "--from="+from)

	delBalAfter := bank.MustGetBalance(t, delAddr, utils.BaseDenom, 0)
	fromBalAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	coinSpent := tx.MustGetDeductFeeAmount(t)
	stakedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	sharedAmount := stakedAmount.Float()

	require.Equal(t, delBalBefore, delBalAfter)
	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(stakedAmount), fromBalAfter)

	valAddr := tx.MustGetEventAttributeValue(t, "create_validator", "validator")
	actualStakedAmount := mustGetStakedAmount(t, tx)

	require.NotEmpty(t, valAddr)
	require.False(t, actualStakedAmount.IsZero())
	require.Equal(t, stakedAmount, actualStakedAmount)

	val := MustGetValidator(t, valAddr)

	require.Equal(t, valPk.Type, val.ConsensusPubkey.Type)
	require.Equal(t, valPk.Key, val.ConsensusPubkey.Key)
	require.Equal(t, commissionRate, val.Commission.CommissionRates.Rate.Float64())
	require.Equal(t, commissionMaxRate, val.Commission.CommissionRates.MaxRate.Float64())
	require.Equal(t, commissionMaxChangeRate, val.Commission.CommissionRates.MaxChangeRate.Float64())
	require.Equal(t, minSelfDelegation, val.MinSelfDelegation)
	require.False(t, val.Jailed)
	require.Equal(t, stakedAmount, val.Tokens)
	require.Equal(t, sharedAmount.String(), val.DelegatorShares.String())

	del := MustGetDelegation(t, delAddr, valAddr)

	require.Equal(t, stakedAmount, del.Balance.Amount)
	require.Equal(t, sharedAmount.String(), del.Delegation.Shares.String())

	return val
}

func MustDelegateForOther(t testutil.TestingT, valAddr string, amount string, from string, delAddr string) {
	valBefore := MustGetValidator(t, valAddr)
	delBalBefore := bank.MustGetBalance(t, delAddr, utils.BaseDenom, 0)
	fromBalBefore := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.MaxBlockTime)
	defer cancel()

	tx := txcmd.MustExecTx(t, ctx, "staking", "delegate-for-other", delAddr, valAddr, amount, "--from="+from)

	valAfter := MustGetValidator(t, valAddr)
	delBalAfter := bank.MustGetBalance(t, delAddr, utils.BaseDenom, 0)
	fromBalAfter := bank.MustGetBalance(t, from, utils.BaseDenom, 0)

	coinSpent := tx.MustGetDeductFeeAmount(t)
	delegatedAmount := testutil.MustGetBaseDenomAmount(t, amount)
	slashedAmount := mustGetSlashedAmount(t, valBefore, valAfter)
	sharedAmount := valBefore.DelegatorShares.Mul(delegatedAmount.DivFloat(valBefore.Tokens))

	require.Equal(t, delBalBefore, delBalAfter)
	require.Equal(t, fromBalBefore.Sub(coinSpent).Sub(delegatedAmount), fromBalAfter)
	require.Equal(t, valBefore.Tokens.Sub(slashedAmount).Add(delegatedAmount), valAfter.Tokens)
	require.Equal(t, valBefore.DelegatorShares.Add(sharedAmount).String(), valAfter.DelegatorShares.String())

	del := MustGetDelegation(t, delAddr, valAddr)

	require.Equal(t, delegatedAmount, del.Balance.Amount)
	require.Equal(t, sharedAmount.String(), del.Delegation.Shares.String())
}

type Pool struct {
	BondedTokens    testutil.Int `json:"bonded_tokens"`
	NotBondedTokens testutil.Int `json:"not_bonded_tokens"`
}

func MustGetStakingPool(t testutil.TestingT) Pool {
	var pool Pool
	cmd.MustQuery(t, &pool, "staking", "pool")
	return pool
}
