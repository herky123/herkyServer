package iface

import "net"

// IConnection 用于处理单个连接
type IConnection interface {
	// Start 启动连接
	Start()
	// Stop 停止
	Stop()
	// GetTCPConnection 获取当前连接信息
	GetTCPConnection() *net.TCPConn
	// GetConnID 获取当前连接Id
	GetConnID() uint32
	// RemoteAddr 获取远程客户端地址信息
	RemoteAddr() net.Addr
}

// HandFunc 用于处理各个连接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
