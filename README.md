# Message Board

## Requirements

- Docker
- cURL

## Check the volume path

Check and fix the volume path in docker-compose.yaml.
Create the directory for the volume path if it does not exist.
```
mkdir -p ~/Desktop/Donat/id/redis-data
```

## Open a terminal in the directory containing docker-compose.yaml

## Build the application

```
go mod tidy
go fmt ./...
docker compose build
```

## Start the application

```
docker compose up -d
```

## Check the services

```
docker compose ps
```

## Check Redis

```
docker exec redis redis-cli PING
```

## Check message-board

```
BASE_URL="http://localhost:8080"

# check health
curl -i "$BASE_URL/healthz"

# create a message and get the message ID
MESSAGE_RESPONSE=$(curl -sS -X POST "$BASE_URL/messages" -H 'Content-Type: application/json' -d '{"title":"title1","body":"body1"}')
echo "MESSAGE_RESPONSE: $MESSAGE_RESPONSE"
MESSAGE_ID=$(printf '%s' "$MESSAGE_RESPONSE" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
echo "MESSAGE_ID: $MESSAGE_ID"

# list messages
curl -i -sS "$BASE_URL/messages"

# create a reply
curl -i -sS -X POST "$BASE_URL/messages/$MESSAGE_ID/replies" -H 'Content-Type: application/json' -d '{"reply":"reply1"}'

# list replies
curl -i -sS "$BASE_URL/messages/$MESSAGE_ID/replies"
```

## Check the message-board logs

```
docker logs -f message-board
```

## Clear the Redis database

```
docker exec redis redis-cli FLUSHDB
```

## Stop the application

```
docker compose down
docker container ls
```

## Further development directions

- Handle non-exisiting message ID in CreateReply and ListReplies and return HTTP 404, save message IDs in a Redis set.
- Add unit tests, mock Redis using GoMock, Testify, or Mockery.
- Test HTTP handlers using http/httptest.
- Add integration tests.
- Validate input JSON.
- Handle timeout of Redis calls. 
- Implement graceful shut down.
- Add panic recovery.
- Version the Rest API in the URL path.
- Add pagination.
- Add Swagger/OpenAPI documentation.
- Upgrade the HTTP connections to WebSocket
- Add authentication and authorization.
- Add user and timestamp to message and reply.
- Secure Redis access.
- Add HTTPS, with or without an API Gateway.
- Set CORS headers.
- Add rate limiting.
- Add Structured JSON logging.
- Add correlation IDs.
- Add metrics endpoint and metrics.
- Check Redis connection pool settings.
- Make HTTP server and Redis configurable through environment variables.

