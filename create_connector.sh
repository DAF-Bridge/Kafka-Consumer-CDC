#!/bin/sh

# Wait for Kafka and PostgreSQL to be available
echo "Waiting for Zookeeper to be available..."
while ! nc -z zookeeper 2181; do sleep 1; done

echo "Waiting for Kafka to be available..."
while ! nc -z kafka 9092; do sleep 1; done

echo "Waiting for PostgreSQL to be available..."
while ! nc -z postgres 5432; do sleep 1; done

echo "Waiting for Debezium to be available..."
while ! nc -z debezium 8083; do sleep 1; done

# Create the Debezium connector
echo "Creating the Debezium PostgreSQL connector..."
curl -X POST http://localhost:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "postgres-cdc-connector",
    "config": {
      "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
      "database.hostname": "postgres",
      "database.port": "5432",
      "database.user": "postgres",
      "database.password": "postgres",
      "database.dbname": "talents_atmos",
      "database.server.name": "cdc",
      "slot.name": "debezium_slot",
      "plugin.name": "pgoutput",
      "table.include.list": "public.events, public.org_open_jobs, public.organizations",
      "topic.prefix": "cdc",
      "database.history.kafka.bootstrap.servers": "localhost:9092",
      "database.history.kafka.topic": "schema-changes.talents_atmos"
    }
  }'

# Wait for the connector to be fully created
echo "Waiting for the connector to be fully created..."
sleep 5

# Check the connector status
echo "Checking the connector status..."
connector_status=$(curl -s http://localhost:8083/connectors/postgres-cdc-connector/status)
echo "Connector Status: $connector_status"

# Check if the connector is running
if echo "$connector_status" | grep '"RUNNING"'; then
    echo "Connector is running successfully."
else
    echo "Connector is not running. Please check the logs for errors."
    exit 1
fi

# Check the connector group
# echo "Checking the connector group..."
# connector_group=$(curl -s http://localhost:8083/connectors/postgres-cdc-connector/tasks)
# echo "Connector Group: $connector_group"

# # Check if the Kafka topic is created
# topic_name="cdc.public.events"
# echo "Checking if Kafka topic $topic_name exists..."
# topic_exists=$(kafka-topics --bootstrap-server localhost:9092 --list | grep -w "$topic_name")
# if [ "$topic_exists" == "$topic_name" ]; then
#     echo "Topic $topic_name exists."
# else
#     echo "Topic $topic_name does not exist."
#     exit 1
# fi

# # Check if data is being consumed (optional)
# echo "Checking if data is being consumed from Kafka..."
# consumer_status=$(curl -s http://debezium:8083/connectors/postgres-cdc-connector/status | jq '.connector.taskStates[0].workerId')
# if [ -z "$consumer_status" ]; then
#     echo "No consumers found for the connector."
# else
#     echo "Consumer found: $consumer_status"
# fi

# echo "Connector setup complete and data is flowing through Kafka."