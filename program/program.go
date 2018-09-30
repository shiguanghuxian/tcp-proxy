package program

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/shiguanghuxian/tcp-proxy/config"
)

// Program 程序实体
type Program struct {
	cfg     *config.Config
	listens []net.Listener // 所有代理列表
	lock    *sync.Mutex    // 获取后端被代理ip时加锁
	sconns  *sync.Map      // 客户端连接列表
	dconns  *sync.Map      // 和后端被代理ip建立的连接列表
}

// New 创建服务对象
func New(cfg *config.Config) *Program {
	return &Program{
		cfg:     cfg,
		listens: make([]net.Listener, 0),
		lock:    new(sync.Mutex),
		sconns:  new(sync.Map),
		dconns:  new(sync.Map),
	}
}

// Start 程序运行
func (p *Program) Start() error {
	if p.cfg == nil {
		return errors.New("配置文件不能为nil")
	}
	defer func() {
		// 打印监听端口映射情况
		var mapProxy string
		for _, v := range p.cfg.Proxys {
			// str := fmt.Sprintf(`%s:\n`, v.Name)
			str := ""
			for _, v1 := range v.Reverse {
				str += fmt.Sprintf("\t%s -> %s\n", v.Listen, v1)
			}
			str = fmt.Sprintf("%s:\n%s\n", v.Name, str)
			mapProxy += str
		}
		fmt.Println(mapProxy)
	}()

	// 启动tcp监听
	err := p.runServer()
	if err != nil {
		return err
	}

	return nil
}

// Stop 停止运行
func (p *Program) Stop() error {
	// 关闭所有建立的连接-客户端
	p.sconns.Range(func(k, v interface{}) bool {
		if v, ok := k.(net.Conn); ok == true {
			err := v.Close()
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})
	// 关闭所有建立的连接-服务端
	p.dconns.Range(func(k, v interface{}) bool {
		if v, ok := k.(net.Conn); ok == true {
			err := v.Close()
			if err != nil {
				log.Println(err)
			}
		}
		return true
	})

	// 关闭所有监听
	for _, v := range p.listens {
		err := v.Close()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
