SERVER="localhost"
PORT="3000"
NUM_REQUESTS=5
DELAY=2  

make_request() {
    RESPONSE=$(curl -s http://$SERVER:$PORT)
    echo "Response from server: $RESPONSE"
}

for ((i=0; i<=NUM_REQUESTS; i++)) do
    echo "Making request #$i to $SERVER:$PORT"
    make_request &
    # sleep $DELAY
done

wait

echo "Finished making $NUM_REQUESTS concurrent requests to $SERVER:$PORT"
