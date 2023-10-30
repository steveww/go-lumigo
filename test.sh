#!/bin/sh

URL="http://localhost:8080"
KEYS="fee fi fo fum"
DATA="tracethisdata"

clear
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
