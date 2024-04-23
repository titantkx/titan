# Upgrade chain

Every change to the chain state is a breaking change. This means that the chain state is not compatible with the previous version of the chain.
It will cause consensus failure if the chain is upgraded without proper migration. Even changed in verify transaction logic will cause consensus failure.

## Versioning

In cosmos base chain, we have version of application and version for each module.

### App version

Version of application will be store in `app/upgrades/[version]/keys.go`. We declare constant `UpgradeName` for each version.

App version will use format `v[major].[minor].[patch]`.
Where major will be increase when we add new module.
Minor will be increase when we change the logic of existing module.
Patch will increase when there is a patch without changing the logic of existing modules.
**NOTE: But we will use format `v[major]_[minor]_[patch]` for every where in code base because cosmos proposal only accept `_` character**

In `app.go` we declare method `setupUpgradeHandlers` to clarify what must todo when upgrade to specific version.
We use `app.UpgradeKeeper.SetUpgradeHandler` to register upgrade handler for each version.

  ```go
  app.UpgradeKeeper.SetUpgradeHandler(
    v1.UpgradeName,
    func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
      return app.mm.RunMigrations(ctx, app.configurator, vm)
    },
  )
  ```

  Or maybe when add new module, we need call `InitGenesis` method of new module (**Only necessary if module is not init from default genesis**) . Example:

  ```go
  app.UpgradeKeeper.SetUpgradeHandler(
    v2.UpgradeName,
    func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

        vm[leaderboardmoduletypes.ModuleName] = leaderboardmodulemigrationscv2types.ConsensusVersion
        genesis, err := leaderboardmodulemigrationscv2.ComputeInitGenesis(ctx, app.CheckersKeeper)
        if err != nil {
            return vm, err
        }
        gen, err := app.appCodec.MarshalJSON(genesis)
        if err != nil {
            return vm, err
        }
        app.mm.Modules[leaderboardmoduletypes.ModuleName].InitGenesis(
            ctx,
            app.appCodec,
            gen)
            
        return app.mm.RunMigrations(ctx, app.configurator, vm)
    },
  )
  ```

`setupUpgradeHandlers` need to be call in `New` method in `app.go`.

An upgrade cannot be process if it do not register in `setupUpgradeHandlers`.

### Module version

Each module version will be understand as consensus version, that is **interger number**, count from 1.
Module version define will be store in `x/[module]/migrations/cv[version]/types/keys.go` where we declare constant `ConsensusVersion` for each version.

Module version must increase by 1 for each upgrade.

Module `github.com/cosmos/cosmos-sdk/x/upgrade` will store version map for every module in database.
So When a module upgrade, it return it own current module version to upgrade module, upgrade module will know what module need to be upgrade.

In `x/[module]/module.go` method `ConsensusVersion` will return current module version.

In `RegisterServices` method we need to use `cfg.RegisterMigration` to register migration for each module version.
  
  ```go
  if err := cfg.RegisterMigration(types.ModuleName, cv1Types.ConsensusVersion, func(ctx sdk.Context) error {
    return migrateV1toV2(ctx, k)
  }); err != nil {
    panic(fmt.Errorf("failed to register migration for %s module: %w", types.ModuleName, err))
  }
  ```

After regis, upgrade module will call this migration function when upgrade module version from old version (`cv1Types.ConsensusVersion`) to old version +1.

## Cosmos upgrade mechanism

1. When some one want to upgrade chain, they deposit some amount of asset to create a proposal to upgrade chain at specific block height .
2. Every proposal have voting period, in voting period, every one can deposit some asset to vote yes or no for this proposal.
3. If a proposal is passed, the chain will halt at block height, every node stop produce new block.
4. Validator run new binary of node with new version of application and new version of module. After migrate process, node will resume produce new block.

## Upgrade steps

### 1. Implement new logic

### 2. Implement migrate process

### 3. Declare new version and regis upgrade handler for module

### 4. Declare new version and regis upgrade handler for app

### 5. Build and publish new binary file for new version

### 6. Create proposal to upgrade chain at specific block height

example:

```shell
titand tx gov submit-legacy-proposal software-upgrade v1tov1_1 --title "v1tov1_1" --description "test" --home ./local_test_data/.titan_val1 --from test --keyring-backend test --chain-id titan_90000-1 --upgrade-height 6000 --deposit 1000titan --upgrade-info '{"binaries":{"any":"file:///Users/mac/Data/Codes/go/tokenize/titan/build/titand_v1_1?checksum=sha256:a597538c45cb1a599ddeb0e0b69c885987abb0222bf16ab78c9b0b2ad2f0ccf5"}}'
```
