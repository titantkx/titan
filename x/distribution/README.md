# `x/distribution`

## Abstract

Extended distribution module of Cosmos-SDK. Add more methods for the keeper and override `AllocateTokens` logic.

## Override `AllocateTokens` logic

In Original distribution module, `AllocateTokens` will allocate tokens to all validators by formula:

```go
totalReward := feePool * (1 - communityTax)

validatorReward := (validatorPower/totalPreviousPower) * totalReward
```

In new logic, we get more rewards from global pool `validator_reward_collector` and distribute them to all validators by formula:

```go
totalReward := feePool * (1 - communityTax) + validatorRewardPool

validatorReward := (validatorPower/totalPreviousPower) * totalReward
```

## New methods

### `FundCommunityPoolFromModule`

Because we remove the burning token mechanism, so all burning tokens will be sent to the community pool.
