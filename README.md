# Chat Server

Chat is a TCP server that allow (netcat "nc") clients to communicate to each other.

Run this on your terminal:
```
  $ go build chat.go
  $ ./chat
```

On a separate instance of your terminal, dial the chat server:
```
  $ nc localhost 8000
```

Run another instance and allow clients to communicate to each other.
