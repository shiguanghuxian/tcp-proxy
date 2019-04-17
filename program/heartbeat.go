package program

import (
	"log"
	"net"
	"time"
)

// 心跳，防止服务端断开连接

// MysqlPing mysql ping 防止mysql服务被代理过，代理服务长时间未操作连接断开
func (p *Program) MysqlPing(conn net.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				// 构建一个ping包
				body := []byte{0x0E}
				body = append(body, []byte("ping")...)
				log.Println("发送mysql ping")
				_, err := conn.Write(body)
				if err != nil {
					log.Println("ping mysql err: ", err)
					return
				}
			}
		}
	}()
}
