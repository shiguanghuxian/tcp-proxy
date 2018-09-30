package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/shiguanghuxian/tcp-proxy/config"
	"github.com/shiguanghuxian/tcp-proxy/program"
)

var f string // 配置文件路径

func main() {
	flag.StringVar(&f, "f", "./config/cfg.yaml", "配置文件路径")
	flag.Parse()
	// 日志输出文件和代码行
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	// 初始化配置文件
	cfg, err := config.NewConfig(f)
	if err != nil {
		log.Fatalln("读取配置文件错误：", err)
	}
	// 创建程序实例
	program := program.New(cfg)
	// 启动
	err = program.Start()
	if err != nil {
		log.Fatalln("程序启动错误：", err)
	}
	log.Println("服务启动成功")

	// 监听退出
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	select {
	case <-c:
		err = program.Stop()
		if err != nil {
			log.Fatalln(err)
		}
	}
}
