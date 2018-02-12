# Chat Server

Chat is a TCP server that allows (netcat "`nc`") clients to communicate to each other.

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

See it in [action](https://lh6.googleusercontent.com/xeW-4ri7e7t_iNY54Q3EVfnhKfMNp4S7FQyJi3-EnYT72WwHO3dbyuYgtgBTFSMY-0klKaDVINt-qh1398tC=w2880-h1472).
