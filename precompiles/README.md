# Precompiles contract

## Generate abi from sol file

1. Install `solc` if you haven't already.
2. install version `0.8.16`

    ```shell
    solc-select install 0.8.16
    ```

3. Set the version to `0.8.16`

    ```shell
    solc-select use 0.8.16
    ```

4. Generate the abi file

    ```shell
    solc --abi [xxx].sol -o .
    ```
