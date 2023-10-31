#!/bin/sh

URL="http://localhost:8080"
KEYS="fee fi fo fum"
DATA="tracethisdata"

usage() {
    echo ""
    echo "Usage:"
    echo "test.sh [-l]"
    echo "-l - loop continuously"
    exit 1
}

unset LOOP
while getopts 'l' option
do
    case "$option" in
        l)
            LOOP=1
            ;;
        *)
            usage
            ;;
    esac
done

clear
while true
do
    echo "Home"
    echo "===="
    curl -w '%{http_code}\n' "${URL}/"

    echo "Adding keys"
    echo "==========="
    for K in $KEYS
    do
        echo "Adding $K"
        curl -w '%{http_code}\n' -X PUT "${URL}/add/${K}/${DATA}-${K}"
    done
    sleep 2

    echo ""
    echo "List keys"
    echo "========="
    curl -w '%{http_code}\n' "${URL}/list"
    sleep 2

    echo ""
    echo "Getting keys"
    echo "============"
    for K in $KEYS
    do
        curl -w '%{http_code}\n' "${URL}/fetch/${K}"
        VAL=$(curl -s "${URL}/fetch/${K}")
        echo "$K = $VAL"
    done
    sleep 2

    echo ""
    echo "Deleting keys"
    echo "============="
    for K in $KEYS
    do
        /bin/echo -n "$K "
        curl -w '%{http_code}\n' -X DEL "${URL}/del/${K}"
    done

    echo ""
    echo "Error"
    echo "====="
    curl -w '%{http_code}\n' "${URL}/add/ping/pong"
    curl -w '%{http_code}\n' "${URL}/fetch/false"

    if [ -z "$LOOP" ]
    then
        break
    fi
    sleep 2
done
