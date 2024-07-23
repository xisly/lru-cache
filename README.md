# LRU-cache 

The following package implements an LRU cache with ttl as an http service

to run:
```bash
go mod download
go run cmd/lru-cache/main.go -log-level="INFO" -server-host-port="localhost:8080" -cache-size=10 -default-cache-ttl=1m
```

to run in docker:
```bash
docker build -t lru-cache .
docker run --env-file .env -p 8080:8080 lru-cache -server-host-port=":8080" -cache-size=100 -log-level="DEBUG"
```
test coverage profile:
```bash
 go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./... && go tool cover -func coverage.out | grep total | awk '{print $3}'
```

There's no background ttl cleanup - to do this cheaply I would use a probabalistic algorithm([kinda like redis does it](https://www.pankajtanwar.in/blog/how-redis-expires-keys-a-deep-dive-into-how-ttl-works-internally-in-redis)), but the problem is that go's ranging over maps isn't truly random. I've found [a really outdated implementation of a random map in go](https://github.com/lukechampine/randmap), maybe I should write a newer one myself.
  
