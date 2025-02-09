package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	Server *Server
}

//create user api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		Server: server,
	}
	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	//加入onlinemap
	this.Server.mapLock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.mapLock.Unlock()
	//广播上线
	this.Server.BroadCast(this, "上线")
}

func (this *User) Offline() {
	this.Server.mapLock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.mapLock.Unlock()

	this.Server.BroadCast(this, "下线")
}

func (this *User) DoMessage(msg string) {
	this.Server.BroadCast(this, msg)
}

//监听User channal
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
