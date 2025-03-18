#!/bin/sh

echo "Waiting for Zookeeper to be available..."
while ! nc -z zookeeper 2181; do sleep 1; done

echo "Waiting for Kafka to be available..."
while ! nc -z kafka 9092; do sleep 1; done

echo "Waiting for PostgreSQL to be available..."
while ! nc -z postgres 5432; do sleep 1; done

echo "Waiting for Debezium to be available..."
while ! nc -z debezium 8083; do sleep 1; done

echo "Waiting for the connector to be fully created..."
sleep 5

echo "Checking the connector status..."
connector_status=$(curl -s http://debezium:8083/connectors/postgres-cdc-connector/status)
echo "Connector Status: $connector_status"

# Check if the connector is running
if echo "$connector_status" | grep '"RUNNING"'; then
    echo "Connector is running successfully."
else
    echo "Connector is not running. Please check the logs for errors."
    exit 1
fi
