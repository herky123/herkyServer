package engine

import (
	"fmt"
	"github.com/herky/herky/iface"
	"net"
)

type Connection struct {
	//当前连接，用于通信的套接字
	Conn *net.TCPConn
	//连接id，SessionID
	ConnId uint32
	//当前连接是否关闭
	isClosed bool
	// 处理请求
	handleAPI iface.HandFunc
	//通知 退出/停止的channel
	ExitBuffChan chan bool
}

// NewConntion 创建一个连接
func NewConntion(conn *net.TCPConn, connID uint32, callback_api iface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnId:       connID,
		isClosed:     false,
		handleAPI:    callback_api,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Printf("Reader Goroutine is running, ConnId = %d\n", c.ConnId)
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit")
	defer c.Stop()
	for {
		//读取数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			c.ExitBuffChan <- true
			return
		}
		//执行handler方法处理业务
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", c.ConnId, " handle is error")
			c.ExitBuffChan <- true
			return
		}
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	//阻塞
	for {
		select {
		//接收到退出消息
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	c.Conn.Close()

	c.ExitBuffChan <- true
	//关闭通道
	close(c.ExitBuffChan)
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnId
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// GetTCPConnection 获取当前连接信息
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}
