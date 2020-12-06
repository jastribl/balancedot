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
  -build="go build -o server cmd/server/main.go" \
  -command="./server" \
  -color=true \
  -graceful-kill=true
```