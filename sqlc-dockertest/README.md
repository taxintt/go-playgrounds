# setup

```
brew install sqlc
```

# check example code

```
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
sqlc generate
```

# run test

```
❯ go test
PASS
ok      github.com/taxintt/go-playgrounds/sqlc-dockertest       45.490s
```

# reference
- https://github.com/tychy/dockertest-sqlc-test-sample/tree/main
- https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html