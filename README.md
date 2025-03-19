# This is a consumer to our Apache Kakfa 
The consumer will consumer the payload from Kafka and process the event and then perform <br />
repository interface with Elasticsearch/Opensearch up to the package of client that you used <br />
within the project of yours. <br />

Do not forget the replace the environment variables according to the .env.example file to run the project.

## Installation Dependencies
After you cloning the project you have to install the package/library that will be using in the project from this command.
```
go mod tidy
```

## Building the project
For building the project you need to run this command
```
go run .\cmd\main.go
```
