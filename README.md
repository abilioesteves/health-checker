# Health Checker

Health Checker is an utility microservice to constantly verify if a target is ok with regards to its dependencies.

Health Checker will periodically call a `/health` endpoint to record its targets health.

The behavior of this endpoint is completely undefined, but requires the following structure as a response:

```
{
    "dependencies": {
        "key": "0|1", # 0 meaning the dependency is down or unreacheable, 1 meaning dependecy is up and reacheable
    }
}
```

# Developer

Launch it fast by running:

```
docker-compose up -d local
```

Or build this application with `go build` and launch it by running:

```
./health-checker start --log-level debug --health-url http://<target>:<port>/health
```

Where `<target>` and `<port>` reflect user defines values for the target base url and port.