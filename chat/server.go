package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms: make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run(){
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	// by default clients have anonymous name, can be changed later
	// commands is the channel that will be used to send commands from client to server
	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	c.readInput()
}

// change client nickname
func (s *server) nick(c *client, args []string){
	c.nick = args[1]
	c.msg(fmt.Sprintf("I will call you %s", c.nick))
}

// join the room
func (s *server) join(c *client, args []string){
	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name: roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = commands

	s.quitCurrentRoom(c)
	c.room = r 
	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

// list all the rooms
func (s *server) listRooms(c *client, args []string){
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms are: %s", strings.Join(rooms, ",")))
}

// send a msg to the current room
func (s *server) msg(c *client, args []string){
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}
	c.room.broadcast(c, c.nick+": " + strings.Join(args[1:len(args)], " "))
}

// quit the chat
func (s *server) quit(c *client, args []string){
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("Sad to see you go")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client){
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}