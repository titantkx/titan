<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
"Miscellaneous" for anything else.

Ref: https://keepachangelog.com/en/1.1.0/
-->

# Changelog

## [Unreleased]

### Bug Fixes

- (tokenfactory) Implement LegacyMsg

## [v3.0.0]

### State Machine Breaking

- (deps) upgrade cometbft to 0.37.6, ibc-go to 7.4.0, gogoproto to 1.4.12
- (ibc) [#81](https://github.com/titantkx/titan/issues/81) Integrate `packetForward` middleware for ibc transfer.
- (ibc) [[#84](https://github.com/titantkx/titan/issues/84)] Integrate `ibcHook` middleware for ibc transfer.
- (tokenfactory) [#88](https://github.com/titantkx/titan/pull/88) Add tokenfactory module.

### Features

- (amino) [#73](https://github.com/titantkx/titan/issues/73) Add regis amino codec of `nftmint` and `validatorreward` module for `gov`,`authz`,`group`.

### Bug Fixes

- (distribution) [#86](https://github.com/titantkx/titan/issues/86) Avoid panic in abci endblocker from `gov` module.

### Miscellaneous

- (Makefile) Update to use fork version of ignite cli.

- Fix simulation tests.

- Add simulation tests for module `validatorreward`.

- (testutil) [#88](https://github.com/titantkx/titan/pull/88) Change function `sample.AccAddress` return type.

- (.golangci.yml) [#88](https://github.com/titantkx/titan/pull/88) Remove rule `allow-leading-space`.

## [v2.0.1](https://github.com/titantkx/titan/releases/tag/v2.0.1)
