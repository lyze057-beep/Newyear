package service

import (
	"5/work/Newyear/user-srv/basic/config"
	__ "5/work/Newyear/user-srv/basic/proto"
	"5/work/Newyear/user-srv/handler/dao"
	"5/work/Newyear/user-srv/handler/model"
	"5/work/Newyear/user-srv/pkg"
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Server struct {
	__.UnimplementedOrderServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) SendSms(_ context.Context, in *__.SendSmsReq) (*__.SendSmsResp, error) {
	code := rand.Intn(9000000) + 1000000
	err := pkg.SendSms(in.Phone, strconv.Itoa(code))
	if err != nil {
		return &__.SendSmsResp{
			Code: int32(code),
			Msg:  "短信发送失败",
		}, nil
	}
	var ctx = context.Background()
	rdb := config.Rdb
	code_key := "sms" + in.Phone
	rdb.Set(ctx, code_key, code, time.Minute*3)
	rabbitmq := pkg.NewRabbitMQSimple("" +
		"code_sms")
	rabbitmq.PublishSimple(strconv.Itoa(code))
	fmt.Println("发送成功！")
	return &__.SendSmsResp{
		Code: int32(code),
		Msg:  "短信发送完毕",
	}, nil
}
func (s *Server) Register(_ context.Context, in *__.RegisterReq) (*__.RegisterResp, error) {
	var ctx = context.Background()
	rdb := config.Rdb
	code_key := "sms" + in.Phone
	result, _ := rdb.Get(ctx, code_key).Result()
	if result != in.Code {
		return nil, fmt.Errorf("验证码错误")
	}
	register := &dao.Register{}
	getUser, err := register.GetUser(in.Username)
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败")
	}
	if getUser != nil {
		return nil, fmt.Errorf("用户已存在")
	}
	hashedPwd, err := pkg.GeneratePassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败")
	}
	newUser := model.User{
		Username: in.Username,
		Password: string(hashedPwd),
		Phone:    in.Phone,
	}
	if err = register.Create(&newUser); err != nil {
		return nil, err
	}
	return &__.RegisterResp{
		Id:  int32(newUser.ID),
		Msg: "注册成功",
	}, nil
}
func (s *Server) Login(_ context.Context, in *__.LoginReq) (*__.LoginResp, error) {
	var ctx = context.Background()
	rdb := config.Rdb
	code_key := "sms" + in.Phone
	result, _ := rdb.Get(ctx, code_key).Result()
	if result != in.Code {
		return nil, fmt.Errorf("验证码错误")
	}
	login := &dao.Register{}
	getUser, err := login.GetUser(in.Username)
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败")
	}
	if getUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	password, err := pkg.ValidatePassword(in.Password, getUser.Password)
	if err != nil {
		return nil, fmt.Errorf("密码对比错误:%v", password)
	}
	token, _ := pkg.GetJwtToken(int64(getUser.ID))
	return &__.LoginResp{
		Id:    int32(getUser.ID),
		Msg:   "登录成功",
		Token: token,
	}, nil
}
