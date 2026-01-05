package main

import (
	_ "5/work/Newyear/user-srv/basic/init"
	__ "5/work/Newyear/user-srv/basic/proto"
	"5/work/Newyear/user-srv/handler/dao"
	"5/work/Newyear/user-srv/handler/service"
	pkg "5/work/Newyear/user-srv/pkg"
	"flag"
	"log"
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 支付回调 HTTP
	go startHTTPNotify()

	s := grpc.NewServer()
	__.RegisterOrderServer(s, &service.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func startHTTPNotify() {
	r := gin.Default()
	r.POST("/alipay/notify", func(c *gin.Context) {
		client, err := pkg.NewAlipayClient()
		if err != nil {
			c.String(200, err.Error())
			return
		}
		if err := c.Request.ParseForm(); err != nil {
			c.String(200, err.Error())
			return
		}
		if err := client.VerifySign(c.Request.PostForm); err != nil {
			c.String(200, err.Error())
			return
		}
		outTradeNo := c.PostForm("out_trade_no")
		tradeNo := c.PostForm("trade_no")
		tradeStatus := c.PostForm("trade_status")
		if tradeStatus == "SUCCESS" || tradeStatus == "FINISHED" {
			o := &dao.OrderDao{}
			if err := o.MarkPaidByOrderNo(outTradeNo, tradeNo); err != nil {
				c.String(200, err.Error())
				return
			}
		}
		c.String(200, "success")
	})
	r.GET("/alipay", func(c *gin.Context) {
		c.String(200, "alipay")
	})
	r.Run(":8080")
}
