#!/bin/bash

# get directory of script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

cd $DIR/../../../

ignite chain init --skip-proto

# kill any existing `titand` processes
process_name="titand"
if pgrep -x "$process_name" > /dev/null
then
    # Get the PID of the process
    pid=$(pgrep -x "$process_name")
    
    # Kill the process
    echo "Killing process $process_name with PID $pid"
    kill $pid
fi

$process_name start --home ./local_test_data/.titan_val1