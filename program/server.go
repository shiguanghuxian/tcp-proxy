package program

import (
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/shiguanghuxian/tcp-proxy/config"
)

// tcp转发实际代码逻辑

// 启动tcp监听服务，监听连接
func (p *Program) runServer() error {
	for _, v := range p.cfg.Proxys {
		v := v
		// 创建一个tcp监听
		lis, err := net.Listen("tcp", v.Listen)
		if err != nil {
			return err
		}
		p.listens = append(p.listens, lis)
		// 等待链接
		go p.waitAccept(v.Name, lis)
	}
	return nil
}

// 等待链接
func (p *Program) waitAccept(name string, lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			log.Println("建立连接错误:%v", err)
			continue
		}
		log.Println("RemoteAddr：", conn.RemoteAddr(), "LocalAddr：", conn.LocalAddr())
		go p.handle(name, conn)
	}
}

// 处理一个链接的tcp转发
func (p *Program) handle(name string, sconn net.Conn) {
	defer sconn.Close()
	// 全局保存连接对象，用于关闭
	p.sconns.Store(sconn, sconn)
	defer p.sconns.Delete(sconn)
	// 设置用不超时
	sconn.SetReadDeadline(time.Time{})
	sconn.SetWriteDeadline(time.Time{})
	// 获取一个后端代理地址
	ip, ok := p.getIP(name)
	if !ok {
		log.Println("获取后端ip地址错误")
		return
	}
	dconn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	// 全局保存连接对象，用于关闭
	p.dconns.Store(dconn, dconn)
	defer p.dconns.Delete(dconn)

	// 当遇到错误时关闭
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		_, err := io.Copy(dconn, sconn)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("连接已经关闭:client")
			} else {
				log.Printf("往%v发送数据失败:%v\n", ip, err)
			}
		}
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		_, err := io.Copy(sconn, dconn)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("连接已经关闭:server")
			} else {
				log.Printf("从%v接收数据失败:%v\n", ip, err)
			}
		}
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	dconn.Close()
}

// 获取一个后端被代理ip
func (p *Program) getIP(name string) (string, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()
	// 找到需要映射的列表
	var proxy *config.Proxy
	for _, v := range p.cfg.Proxys {
		if v.Name == name {
			proxy = v
		}
	}
	if len(proxy.Reverse) < 1 {
		return "", false
	}
	ip := proxy.Reverse[0]
	proxy.Reverse = append(proxy.Reverse[1:], ip)
	log.Println(ip)
	log.Println(proxy.Reverse)
	return ip, true
}
