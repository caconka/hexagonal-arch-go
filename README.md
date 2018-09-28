# hex-arch

Example of hexagonal architecture - DDD in go

#### cURL
```
curl -X POST http://localhost:3000/tickets -H 'Cache-Control: no-cache' \
	-H 'Content-Type: application/json' -d '{
	"creator" : "Ted",
	"title" : "Test",
	"description" : "A test ticket",
	"points": 3
}'
```

```
curl http://localhost:3000/tickets
```
