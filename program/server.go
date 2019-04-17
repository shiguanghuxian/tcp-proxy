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
	defer func() {
		dconn.Close()
		p.dconns.Delete(dconn)
	}()

	// 获取连接类型，如果是mysql加入ping
	oneProxy := p.cfg.GetProxyByName(name)
	if oneProxy != nil {
		if oneProxy.Typ == "mysql" {
			p.MysqlPing(dconn)
		}
	}

	// 当遇到错误时关闭
	exitChan := make(chan bool, 1)
	// 客户端->服务端
	go p.tcpIOCopy(sconn, dconn, "客户端->服务端", exitChan)
	// 服务端->客户端
	go p.tcpIOCopy(dconn, sconn, "服务端->客户端", exitChan)

	<-exitChan

}

// tcp copy
func (p *Program) tcpIOCopy(src, dst net.Conn, direction string, exitChan chan bool) {
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)
	var err error
	var n int

	for {
		n, err = src.Read(buf)
		if n > 0 {
			// 写入读取的字节
			_, err = dst.Write(buf[0:n])
			if err != nil {
				log.Println(direction, "写流错误")
				break
			}
		}

		if err != nil || n == 0 {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			if err != io.EOF {
				log.Println(direction, "读取流错误")
			}
			break
		}
	}
	exitChan <- true
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
