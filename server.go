package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Client struct {
	C    chan string //用于发送数据的管道
	Name string      //用户名
	Addr string      //网络地址
}

//保存在线用户 cliAddr ==> Client
var onlineMap map[string]Client
var message = make(chan string)

//转发消息，只要有消息来了，遍历map，给map每个成员都发送此消息
func Manager() {
	//给map分配空间
	onlineMap = make(map[string]Client)
	for {
		msg := <-message //没有消息前，这里会阻塞
		//遍历map，给map每个成员都发送此消息
		for _, cli := range onlineMap {
			cli.C <- msg
		}
	}
}
func WriteMsgToClient(cli Client, conn net.Conn) {
	for msg := range cli.C { //给当前客户端发送信息
		conn.Write([]byte(msg + "\n"))
	}
}
func MakeMsg(cli Client, msg string) (buf string) {
	buf = "[" + cli.Name + "]" + " : " + msg
	return
}
func HandleConn(conn net.Conn) { //处理用户连接
	defer conn.Close()
	//获取客户端的网络地址
	cliAddr := conn.RemoteAddr().String()
	//创建一个结构体,默认用户名和网络地址一样
	cli := Client{make(chan string), cliAddr, cliAddr}
	//把结构体添加到map
	onlineMap[cliAddr] = cli
	//新开一个协程，专门给当前客户端发送信息
	go WriteMsgToClient(cli, conn)
	//广播某个在线
	//message <- "[" + cli.Addr + "]" + cli.Name + " : login"
	message <- MakeMsg(cli, "login")
	//提示,我是谁
	cli.C <- MakeMsg(cli, "Hello! My friend! \nYou can get help by entering the keyword 'help'!")
	isQuit := make(chan bool)  //对方是否主动退出
	hasData := make(chan bool) //对方是否有数据发送
	//新建一个协程，接收用户发送过来的数据
	go func() {
		buf := make([]byte, 2048)
		var help string = `
USAGE:
	"who" -- Get a list of online members
	"rename newname"  -- Set a new sickname,eg : rename 小明
	Default,chat mode! You can chat with others online members!! 
`
		for {
			n, err := conn.Read(buf)
			if n == 0 { //对方断开，或者出问题
				isQuit <- true
				fmt.Println("conn.Read err=", err)
				return
			}
			msg := string(buf[:n-2]) //通过Windows nc测试，多一个换行
			if len(msg) == 3 && msg == "who" {
				//遍历map，给当前用户发送所有成员
				conn.Write([]byte("  User list :\n"))
				for _, tmp := range onlineMap {
					msg = "[" + tmp.Addr + "]" + ":" + tmp.Name + "\n"
					conn.Write([]byte(msg))
				}
			} else if len(msg) == 4 && msg == "help" {
				conn.Write([]byte(help))
			} else if len(msg) >= 8 && msg[:6] == "rename" {
				//rename|mike
				name := strings.Split(msg, " ")[1]
				cli.Name = name
				onlineMap[cliAddr] = cli
				conn.Write([]byte("rename ok\n"))
			} else {
				//转发此内容
				message <- MakeMsg(cli, msg)
			}
			hasData <- true //代表有数据
		}
	}() //别忘了（）
	for {
		//通过select检测channnel的流动
		select {
		case <-isQuit:
			delete(onlineMap, cliAddr) //当前用户从map移除
			message <- MakeMsg(cli, "login out")
			return
		case <-hasData:
		case <-time.After(300 * time.Second):
			delete(onlineMap, cliAddr) //当前用户从map移除
			message <- MakeMsg(cli, "time out,leave out")
			return
		}
	}
}
func main() {
	//监听
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listener.Close()
	//新开一个协程，转发消息，只要有消息来了，遍历map，给map每个成员都发送此消息
	go Manager()
	//主协程，循环阻塞等待用户连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err=", err)
			continue
		}
		go HandleConn(conn) //处理用户连接
	}
}
