package service

import (
	"5/work/Newyear/user-srv/basic/config"
	__ "5/work/Newyear/user-srv/basic/proto"
	"5/work/Newyear/user-srv/handler/dao"
	"5/work/Newyear/user-srv/handler/model"
	"5/work/Newyear/user-srv/pkg"
	"context"
	"errors"
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
func (s *Server) UpdatePassword(_ context.Context, in *__.UpdatePasswordReq) (*__.UpdatePasswordResp, error) {
	updatedao := &dao.UpdatePassword{}
	bytes, err := pkg.GeneratePassword(in.NewPassword)
	if err != nil {
		return nil, errors.New("密码加密失败:%v" + err.Error())
	}
	hashedPwd := string(bytes)
	err = updatedao.UpdateDao(in.Username, hashedPwd)
	if err != nil {
		return &__.UpdatePasswordResp{
			Code: 400,
			Msg:  "密码修改失败",
		}, nil
	}
	return &__.UpdatePasswordResp{
		Code: 200,
		Msg:  "密码修改成功",
	}, nil
}
func (s *Server) ListUser(_ context.Context, in *__.ListUserReq) (*__.ListUserResp, error) {
	listDao := &dao.Register{}
	getUser, err := listDao.GetUser(in.Username)
	if err != nil {
		return nil, fmt.Errorf("查询用户失败:" + err.Error())
	}
	if getUser == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	userInfo := &__.UserInfo{
		Username: getUser.Username,
		Phone:    getUser.Phone,
	}
	return &__.ListUserResp{
		Code: 200,
		Msg:  "查询成功",
		Data: userInfo,
	}, nil
}
func (s *Server) ListGet(_ context.Context, in *__.ListGetReq) (*__.ListGetResp, error) {
	listDao := &dao.ListDao{}
	listUser, err := listDao.ListUser(config.DB)
	if err != nil {
		return nil, fmt.Errorf("查询失败:" + err.Error())
	}
	if len(listUser) == 0 {
		return nil, fmt.Errorf("数据不存在")
	}
	protoUsers := make([]*__.ListUser, 0, len(listUser))
	for _, user := range listUser {
		protoUsers = append(protoUsers, &__.ListUser{
			Username: user.Username,
			Phone:    user.Phone,
			Id:       int32(user.ID),
		})
	}
	return &__.ListGetResp{
		Code: 0,
		Msg:  "查询成功",
		Data: protoUsers,
	}, nil
}
func (s *Server) Product(_ context.Context, in *__.ProductReq) (*__.ProductResp, error) {
	productDao := &dao.Product{}
	pro := &model.Product{
		ProductName:  in.ProductName,
		ProductPrice: int(in.ProductPrice),
		ProductNum:   int(in.ProductNum),
		Status:       int(in.Status),
	}
	if err := productDao.CreateProduct(pro); err != nil {
		return nil, fmt.Errorf("商品添加失败")
	}
	return &__.ProductResp{
		ProductId: int32(pro.ID),
		Code:      200,
		Msg:       "商品添加成功",
	}, nil
}
func (s *Server) ListProduct(_ context.Context, in *__.ListProductReq) (*__.ListProductResp, error) {
	listDao := &dao.Product{}
	get, err := listDao.Get(uint(in.Id))
	if err != nil {
		return nil, fmt.Errorf("id查询失败")
	}
	if get == nil {
		return nil, fmt.Errorf("数据库不存在该商品")
	}
	getProduct := &__.GetProduct{
		ProductName:  get.ProductName,
		ProductPrice: int32(get.ProductPrice),
		ProductNum:   int32(get.ProductNum),
		Status:       int32(get.Status),
		Id:           int32(get.ID),
	}
	return &__.ListProductResp{
		Code: 200,
		Msg:  "查询成功",
		Data: getProduct,
	}, nil
}
func (s *Server) GetProduct(_ context.Context, in *__.GetProductReq) (*__.GetProductResp, error) {
	getDao := &dao.Product{}
	list, err := getDao.List()
	if err != nil {
		return nil, fmt.Errorf("查询失败:" + err.Error())
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("数据库为空")
	}
	Productlist := make([]*__.Product, 0, len(list))
	for _, product := range list {
		Productlist = append(Productlist, &__.Product{
			ProductName:  product.ProductName,
			ProductPrice: int32(product.ProductPrice),
			ProductNum:   int32(product.ProductNum),
			Status:       int32(product.Status),
			Id:           int32(product.ID),
		})
	}
	return &__.GetProductResp{
		Code: 200,
		Msg:  "查询成功",
		Data: Productlist,
	}, nil
}
