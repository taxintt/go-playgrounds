# 0. run docker container

```
❯ docker compose up -d 
```

# 1. set data

```
 ❯ curl -X POST -H "Content-Type: application/json" -d '{"num": 123}' http://localhost:8080/set
{"message":"value inserted"}
```

# 2. refer data

```
❯ curl -X GET http://localhost:8080/get?num=123
{"message":"inserted"}
❯ curl -X GET http://localhost:8080/get?num=456
value not found
```