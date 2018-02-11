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

![Chat in action](https://lh5.googleusercontent.com/Gdr1DtOTQvuiyBO_lT_gew5NDg7zxHmQtzSwbqrJ5quLG8dVPOmPE0iVAfvfSOgkyVmRd6iQ58S1NChjoPb4=w2880-h1472)
