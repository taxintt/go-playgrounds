# 1. run server
```bash
❯ go run test.go 
awaiting signal
working with &{st:test}
working with &{st:test}
working with &{st:test}

...
```

# 2. use pprof 
```bash
❯ go tool pprof -seconds 60 http://localhost:8080/debug/pprof/profile
Fetching profile over HTTP from http://localhost:8080/debug/pprof/profile?seconds=60
Please wait... (1m0s)
Saved profile in /Users/taxin/pprof/pprof.samples.cpu.001.pb.gz
Type: cpu
Time: Oct 29, 2023 at 2:06pm (JST)
Duration: 60s, Total samples = 20ms (0.033%)
Entering interactive mode (type "help" for commands, "o" for options)
```
