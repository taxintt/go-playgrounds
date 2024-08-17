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

# start sample application
```
docker compose up 
```

# create User
```
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Kuro","age":10}' 'localhost:8080/user'
user was created: User(id=1, age=10, name=Kuro)%                             
```

# get all users info
```
curl 'localhost:8080/user/1'
user returned: User(id=1, age=10, name=Kuro)%                                
```

# References
- [Quick Introduction | ent](https://entgo.io/ja/docs/getting-started/)
- [ent を利用している project の migration に atlas を使ってみる - Techtouch Developers Blog](https://tech.techtouch.jp/entry/ent-atlas-migration)