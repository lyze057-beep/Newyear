package main

import (
	__ "5/work/Newyear/user-srv/basic/proto"
	"5/work/Newyear/user-srv/handler/dao"
	"5/work/Newyear/user-srv/handler/service"
	"5/work/Newyear/user-srv/pkg"
	"flag"
	"log"
	"net"
	"strconv"

	_ "5/work/Newyear/user-srv/basic/init"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.

func main() {
	flag.Parse()
	go startHTTPNotify()
	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		rabbitmq := pkg.NewRabbitMQSimple("" +
			"code_sms")
		rabbitmq.ConsumeSimple()
	}()
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
		err = c.Request.ParseForm()
		if err := client.VerifySign(c.Request.PostForm); err != nil {
			c.String(200, err.Error())
			return
		}
		orderNo := c.PostForm("orderNo")
		tradeNo := c.PostForm("tradeNo")
		payStatus := c.PostForm("payStatus")
		if payStatus == "SUCCESS" || payStatus == "FAIL" {
			o := &dao.OrderDao{}
			if err := o.MarkPaid(orderNo, tradeNo); err != nil {
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
