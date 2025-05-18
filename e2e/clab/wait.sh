#!/bin/bash
MAX_TIME=300
INTERVAL=5
END_TIME=$((SECONDS + MAX_TIME))

while [ $SECONDS -lt $END_TIME ]; do
    echo "waiting for ssh..."
    # ensure srl is up and running and displays the normal banner (is past the initial boot stuff
    # basically where it is only basic cli)
    OUT=$(sshpass -p "NokiaSrl1!" \
      ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no -l admin 172.20.20.16 quit  2>&1)

    echo "read: $OUT"

    if echo "$OUT" | grep -q "Welcome to Nokia SR Linux"; then
        echo "ssh available..."
        break
    fi

    sleep $INTERVAL
done

while [ $SECONDS -lt $END_TIME ]; do
    echo "waiting for netconf..."
    # same thing for netconf port just to be safe!
    OUT=$(sshpass -p "NokiaSrl1!" \
      ssh -o ConnectTimeout=5 -o StrictHostKeyChecking=no -l admin -p 830 172.20.20.16 quit  2>&1)

    echo "read: $OUT"

    if echo "$OUT" | grep -q "Welcome to Nokia SR Linux"; then
        echo "netconf available..."
        exit 0
    fi

    sleep $INTERVAL
done

exit 1
