#!/usr/bin/env bash


for i in 70
do

for j in 50
do
echo "{\"name\":\"ROTATE\",\"arguments\":[\"$i\", \"$j\"]}" | nc -4u -w1 10.10.10.20 10001
sleep 0.5
done

sleep 0.5

done;


for i in 0
do

for j in 100
do
echo "{\"name\":\"ROTATE\",\"arguments\":[\"$i\", \"$j\"]}" | nc -4u -w1 10.10.10.20 10001
sleep 0.5
done

sleep 0.5

done;





#
#
#echo -n '{"name":"ROTATE","arguments":["0", "100"]}' | nc -4u -w1 10.10.10.10 10001 && \
#sleep 2 && \
#echo -n '{"name":"ROTATE","arguments":["10", "100"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["20", "100"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["30", "100"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["40", "90"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["50", "80"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["50", "70"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["50", "60"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["60", "50"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["70", "50"]}' | nc -4u -w1 10.10.10.10 10001 && \
#echo -n '{"name":"ROTATE","arguments":["80", "80"]}' | nc -4u -w1 10.10.10.10 10001