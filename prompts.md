# Message Board prompts

## Prompt 1

Find the latest available Docker container version of Redis.
Create a Docker Compose YAML file for this version of Redis and name the file docker-compose.yaml.
Redis should use append-only persistence.
Use the following directory for persistent data: ~/Desktop/Donat/id/redis-data.
Add a comment to adapt this directory to the user's system.
Generate commands to test the Redis container.

## Prompt 2

Create the following files: /cmd/main.go, /internal/web/rest.go, /internal/repository/redis.go, and /internal/repository/repository.go.
The repository.go file defines a Repository interface with a Ping() method that returns a string and an error.
The redis.go file implements the Repository interface with a Ping() method that connects to the Redis instance, sends a ping, returns the response or an error, and closes the connection.
The rest.go file defines a Rest struct, accepts the Repository interface, and it has a healthz() HTTP handler method that returns HTTP 200 when the Ping() method returns PONG, otherwise it returns HTTP 500. It should also log the error.
The main.go file reads the REDIS_URL environment variable, creates a Redis repository instance, creates the REST server with the repository injected into it, starts an HTTP server, and binds the health handler to GET /healthz.
Do not use any frameworks.
The HTTP server will be extended later with additional methods, so use an explicit struct for the server to make it configurable.
Generate a DOCKERFILE with multistage build. Use scratch as the base image for the second stage. Expose the application on port 8080.
Extend the previously generated docker-compose.yaml with the new service and name it message-board. It should communicate with the previously defined Redis service.
Add a new bridge network named message-board-network to both services in docker-compose.yaml.
Add REDIS_URL as an environment variable to message-board service in docker-compose.yaml.
Create a .dockerignore file to ignore txt, md, and yaml files.
Use Go version is 1.26.5.

## Manual fixes

- Dockerfile: COPY go.sum was missing.
- redis.go: In Ping(), redis.ParseURL(r.url) wanted to create a client instead of clientOptions.

## Playing with Redis to explore a data structure

```
docker exec redis redis-cli RPUSH messageboard '{"id":"id1","title":"title1","body":"body1"}'
docker exec redis redis-cli RPUSH messageboard '{"id":"id2","title":"title2","body":"body2"}'
docker exec redis redis-cli LRANGE messageboard 0 -1
docker exec redis redis-cli RPUSH messageboard:id1 '{"reply":"reply1"}'
docker exec redis redis-cli RPUSH messageboard:id1 '{"reply":"reply2"}'
docker exec redis redis-cli LRANGE messageboard:id1 0 -1
docker exec redis redis-cli LRANGE messageboard:id2 0 -1
docker exec redis redis-cli FLUSHDB
```

## Prompt 3
Create the following methods in /internal/repository/repository.go and implement them in /internal/repository/redis.go:
- CreateMessage(title, body string) (string, error)
- ListMessages() (string, error)
- CreateReply(messageID, reply string) error
- ListReplies(messageID string) (string, error)
Use the following Redis operations:
- CreateMessage: redis-cli RPUSH messageboard '{"id":"uuid1","title":"title1","body":"body1"}'
- ListMessages: redis-cli LRANGE messageboard 0 -1
- CreateReply: redis-cli RPUSH messageboard:uuid1 '{"reply":"reply1"}'
- ListReplies: redis-cli LRANGE messageboard:uuid1 0 -1
CreateMessage should generate a UUID, create JSON containing id, title, and body, push it to messageboard, and return the UUID.
CreateReply should create JSON containing reply, and push it to messageboard:<messageId>.
ListMessages and ListReplies must return the Redis entries as a JSON array string.
Rename healthz to healthzHandler in rest.go and add these handlers:
- createMessageHandler: POST /messages
- listMessagesHandler: GET /messages
- createReplyHandler: POST /messages/{messageID}/replies
- listRepliesHandler: GET /messages/{messageID}/replies
Each handler must call the corresponding method on the injected repository.Repository interface.
createMessageHandle accepts JSON containing title and body.
createReplyHandler accepts JSON containing reply.
listMessagesHandler and listRepliesHandler return the JSON received from the repository.
Generate cURL commands to test the endpoints.

## Manual fixes

- Fixed create Redis client creation again.
- Refactored JSON decoding in rest.go, it was over-engineered.
- Removed the custom UUID generation and replaced it with github.com/google/uuid.

## Playing with Redis to explore how to store messageIDs in a set

```
docker exec redis redis-cli SADD messageboard-ids "id1"
docker exec redis redis-cli SISMEMBER messageboard-ids "id1"
docker exec redis redis-cli SISMEMBER messageboard-ids "id2"
docker exec redis redis-cli FLUSHDB
```
