## TCP chat in Go

Once the user connects to the server it can use any of the followint commands to interact with it:

- /nick <name> - get a name, otherwise user will stay anonymous.
- /join <name> - join a room, if room doesn't exist, the new room will be created. User can be only in one room at the same time.
- /rooms - show list of available rooms to join.
- /msg <msg> - broadcast message to everyone in a room.
- /quit - disconnects from the chat server.

*To run*

Server:
```$ go build .```
```$ ./chat```

Clients:
``` telnet localhost 8888```
