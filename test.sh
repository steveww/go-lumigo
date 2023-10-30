#!/bin/sh

URL="http://localhost:8080"
KEYS="fee fi fo fum"
DATA="tracethisdata"

clear
echo "Home"
curl "${URL}/"

echo "Adding keys"
for K in $KEYS
do
    echo "Adding $K"
    curl -X PUT "${URL}/add/${K}/${DATA}${K}"
done
sleep 2

echo ""
echo "List keys"
curl "${URL}/list"
sleep 2

echo ""
echo "Getting keys"
for K in $KEYS
do
    VAL=$(curl -s "${URL}/fetch/${K}")
    echo "$K = $VAL"
done
sleep 2

echo ""
echo "Deleting keys"
for K in $KEYS
do
    /bin/echo -n "$K "
    curl -X DEL "${URL}/del/${K}"
done
