# `x/validatorreward`

## Abstract

This module supports a reward pool for validators. The reward is distributed to all contribute validators using the rate per year that is configured.

## State

### Params

- `rate`: percentage of staking tokens that are distributed to validators per year
- `authority`: the address that allows to configure the params

### `LastDistributeTime`: time of the previous block

## Begin-Block

At every beginning block, the module will calculate the total VotingPower of all signed validators in the previous block.

The total distributed reward in the current block will be calculated by the formula:

```go
totalRewardNeed := totalVotingPower * PowerReduction * rate * ((currentBlockTime - lastDistributeTime) / timePerYear)

totalReward := Min(totalRewardNeed, rewardPool)
```

- `PowerReduction`: The factor that converts Power to amount of staking tokens.
- `timePerYear`: The time of a year, calculated as 365 days, 24 hours per day.

The reward will be sent to global pool `validator_reward_collector` and then distributed to validators in same block by module `distribution`.
