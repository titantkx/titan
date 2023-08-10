package keeper

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Slash a validator for an infraction committed at a known height
// Find the contributing stake at that height and burn the specified slashFactor
// of it, updating unbonding delegations & redelegations appropriately
//
// CONTRACT:
//
//	slashFactor is non-negative
//
// CONTRACT:
//
//	Infraction was committed equal to or less than an unbonding period in the past,
//	so all unbonding delegations and redelegations from that height are stored
//
// CONTRACT:
//
//	Slash will not slash unbonded validators (for the above reason)
//
// CONTRACT:
//
//	Infraction was committed at the current height or at a past height,
//	not at a height in the future
func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) math.Int {
	logger := k.Logger(ctx)

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := k.TokensFromConsensusPower(ctx, power)
	slashAmountDec := sdk.NewDecFromInt(amount).Mul(slashFactor)
	slashAmount := slashAmountDec.TruncateInt()

	// ref https://github.com/cosmos/cosmos-sdk/issues/1348

	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(
			"WARNING: ignored attempt to slash a nonexistent validator; we recommend you investigate immediately",
			"validator", consAddr.String(),
		)
		return sdk.NewInt(0)
	}

	// should not be slashing an unbonded validator
	if validator.IsUnbonded() {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.GetOperator()))
	}

	operatorAddress := validator.GetOperator()

	// call the before-modification hook
	if err := k.Hooks().BeforeValidatorModified(ctx, operatorAddress); err != nil {
		k.Logger(ctx).Error("failed to call before validator modified hook", "error", err)
	}

	// Track remaining slash amount for the validator
	// This will decrease when we slash unbondings and
	// redelegations, as that stake has since unbonded
	remainingSlashAmount := slashAmount

	switch {
	case infractionHeight > ctx.BlockHeight():
		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():
		// Special-case slash at current height for efficiency - we don't need to
		// look through unbonding delegations or redelegations.
		logger.Info(
			"slashing at current height; not scanning unbonding delegations & redelegations",
			"height", infractionHeight,
		)

	case infractionHeight < ctx.BlockHeight():
		// Iterate through unbonding delegations from slashed validator
		unbondingDelegations := k.GetUnbondingDelegationsFromValidator(ctx, operatorAddress)
		for _, unbondingDelegation := range unbondingDelegations {
			amountSlashed := k.SlashUnbondingDelegation(ctx, unbondingDelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}

			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}

		// Iterate through redelegations from slashed source validator
		redelegations := k.GetRedelegationsFromSrcValidator(ctx, operatorAddress)
		for _, redelegation := range redelegations {
			amountSlashed := k.SlashRedelegation(ctx, validator, redelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}

			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}
	}

	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(remainingSlashAmount, validator.Tokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, math.ZeroInt()) // defensive.

	// we need to calculate the *effective* slash fraction for distribution
	if validator.Tokens.IsPositive() {
		effectiveFraction := sdk.NewDecFromInt(tokensToBurn).QuoRoundUp(sdk.NewDecFromInt(validator.Tokens))
		// possible if power has changed
		if effectiveFraction.GT(math.LegacyOneDec()) {
			effectiveFraction = math.LegacyOneDec()
		}
		// call the before-slashed hook
		if err := k.Hooks().BeforeValidatorSlashed(ctx, operatorAddress, effectiveFraction); err != nil {
			k.Logger(ctx).Error("failed to call before validator slashed hook", "error", err)
		}
	}

	// Deduct from validator's bonded tokens and update the validator.
	// Burn the slashed tokens from the pool account and decrease the total supply.
	validator = k.RemoveValidatorTokens(ctx, validator, tokensToBurn)

	switch validator.GetStatus() {
	case types.Bonded:
		// if err := k.burnBondedTokens(ctx, tokensToBurn); err != nil {
		// 	panic(err)
		// }
	case types.Unbonding, types.Unbonded:
		// if err := k.burnNotBondedTokens(ctx, tokensToBurn); err != nil {
		// 	panic(err)
		// }
	default:
		panic("invalid validator status")
	}

	logger.Info(
		"validator slashed by slash factor",
		"validator", validator.GetOperator().String(),
		"slash_factor", slashFactor.String(),
		"burned", tokensToBurn,
	)
	return tokensToBurn
}
