#!/bin/bash

LOG_FILE_PATH="./genlog.log"
I=0

for (( ; ; ))
do
    CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    echo "$CURRENT_TIME Writing a line to $LOG_FILE_PATH..."
    echo "$CURRENT_TIME Log line $I" >> $LOG_FILE_PATH
    I=$(($I+1))
    sleep 1
done
