# setup
```
go install entgo.io/ent/cmd/ent@latest

brew install ariga/tap/atlas
```

```
❯ atlas version
atlas version v0.17.1-8e610d9-canary
https://github.com/ariga/atlas/releases/latest

❯ cat go.mod| grep ent                         
module github.com/taxintt/go-playgrounds/ent-with-atlas
        entgo.io/ent v0.12.5 // indirect
        github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
```

# Create first schema 
```
go run entgo.io/ent/cmd/ent new User
```

# References
- https://tech.techtouch.jp/entry/ent-atlas-migration