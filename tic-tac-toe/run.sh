#!/bin/sh

# Start the Tic-Tac-Toe web application in the background
./tic-tac-toe-app &

# Start the Kafka consumer (foreground)
./kafka-consumer
