package engine

import (
	"errors"
	"fmt"
	"github.com/herky/herky/iface"
	"net"
	"time"
)

type Server struct {
	// ip
	IP string
	//端口
	Port int
	//服务器名
	Name string
	//IP版本 ip4 ip6
	IpVersion string
}

//实现iserver接口

// 开始
func (s *Server) Start() {
	fmt.Printf("###Server engine begin start, ip: %s; port : %d; name : %s; ipVersion: %s\n", s.IP, s.Port,
		s.Name, s.IpVersion)
	//开启一个协程处理
	go func() {
		addr, err := net.ResolveTCPAddr(s.IpVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("addr error : ", err)
			return
		}
		listenner, err := net.ListenTCP(s.IpVersion, addr)
		if err != nil {
			fmt.Println("listen error : ", err)
			return
		}
		fmt.Println("server engine start success, listen start")
		//循环监听
		var cid uint32
		cid = 0
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("accept error :", err)
				continue
			}
			//开启协程处理连接
			dealConn := NewConntion(conn, cid, CallBackToClient)
			cid++
			go dealConn.Start()
		}
	}()
}

// 停止
func (s *Server) Stop() {

}

func (s *Server) Serve() {
	s.Start()
	//阻塞
	for {
		time.Sleep(time.Second * 10)
	}
}

func NewServer(name string) iface.IServer {
	s := &Server{
		Name:      name,
		IpVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      7777,
	}
	return s
}

func CallBackToClient(conn *net.TCPConn, buffer []byte, cnt int) error {
	//回显
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(buffer[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}
