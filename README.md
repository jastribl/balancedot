# BalanceDot

Application for managing expenses with Chase and Splitwise

# Run

Front end
```
cd client && npm run watch
```

Server
```
env $(cat db.env) CompileDaemon \
  -directory=./ \
  -exclude-dir="client" \
  -exclude-dir=".git" \
  -build="go build -o build/server cmd/server/main.go" \
  -command="./build/server" \
  -color=true \
  -graceful-kill=true
```
