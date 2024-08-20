# Interchain tests

## Step

1. Build local image for `titand`

    Run the following command at the root project folder to build the local image for `titand`:

    ```shell
    docker build -t docker.io/titantkx/titand:local .
    ```

2. Run the interchain tests

    Run the following command at the root project folder to run the interchain tests:

    ```shell
    make test-interchain
    ```
