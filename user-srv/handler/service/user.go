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

	"github.com/redis/go-redis/v9"
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
func (s *Server) OrderAdd(_ context.Context, in *__.OrderAddReq) (*__.OrderAddResp, error) {
	orderdao := &dao.OrderDao{}
	var order *model.OrderMain
	var err error
	if in.ProductId > 0 && in.Quantity > 0 {
		order, err = orderdao.CreateOrderWithTx(uint(in.UserId), uint(in.ProductId), int(in.Quantity))
	} else {
		order, err = orderdao.CreateOrderFromCartTx(uint(in.UserId))
	}
	if err != nil {
		return nil, fmt.Errorf("订单添加失败:" + err.Error())
	}
	var ctx = context.Background()
	rdb := config.Rdb
	orderKey := "order:" + order.OrderNo
	_, _ = rdb.HSet(ctx, orderKey, map[string]interface{}{
		"Id":          order.ID,
		"OrderNo":     order.OrderNo,
		"UserId":      order.UserID,
		"TotalAmount": order.TotalAmount,
		"PayStatus":   order.PayStatus,
		"OrderStatus": order.OrderStatus,
		"TradeNo":     order.TradeNo,
	}).Result()
	userOrdersKey := fmt.Sprintf("user:%d:orders", order.UserID)
	_, _ = rdb.ZAdd(ctx, userOrdersKey, redis.Z{Score: float64(time.Now().Unix()), Member: order.OrderNo}).Result()
	amount := fmt.Sprintf("%.2f", order.TotalAmount/10.0)
	alipay, _ := pkg.Alipay("订单支付", order.OrderNo, amount)
	return &__.OrderAddResp{
		OrderId: int32(order.ID),
		Code:    200,
		Msg:     "订单添加成功",
		Data:    alipay,
	}, nil
}
func (s *Server) GetOrder(_ context.Context, in *__.GetOrderReq) (*__.GetOrderResp, error) {
	getdao := &dao.OrderDao{}
	getOrder, err := getdao.GetOrder(int(in.PayStatus))
	if err != nil {
		return nil, fmt.Errorf("数据库查询失败:" + err.Error())
	}
	if getOrder == nil {
		return nil, fmt.Errorf("这条信息不存在")
	}
	get := &__.GetOrder{
		Id:          int32(getOrder.ID),
		OrderNo:     getOrder.OrderNo,
		UserId:      int32(getOrder.UserID),
		TotalAmount: float32(getOrder.TotalAmount),
		PayStatus:   int32(getOrder.PayStatus),
		OrderStatus: int32(getOrder.OrderStatus),
		TradeNo:     getOrder.TradeNo,
	}
	return &__.GetOrderResp{
		Code: 200,
		Msg:  "查询成功",
		Data: get,
	}, nil
}
func (s *Server) CartAdd(_ context.Context, in *__.CartAddReq) (*__.CartAddResp, error) {
	fmt.Printf("收到请求: UserID=%d, ProductID=%d, Quantity=%d\n", in.UserID, in.ProductID, in.Quantity)
	var ctx = context.Background()
	rdb := config.Rdb
	key := fmt.Sprintf("cart:%d", in.UserID)
	field := strconv.Itoa(int(in.ProductID))
	_, _ = rdb.HSet(ctx, key, field, in.Quantity).Result()
	cartDao := &dao.CartDao{}
	if err := cartDao.Add(uint(in.UserID), uint(in.ProductID), int(in.Quantity)); err != nil {
		return nil, fmt.Errorf("购物车添加失败")
	}
	return &__.CartAddResp{
		Code: 200,
		Msg:  "购物车添加成功",
	}, nil
}
func (s *Server) CartList(_ context.Context, in *__.CartListReq) (*__.CartListResp, error) {
	var ctx = context.Background()
	rdb := config.Rdb
	key := fmt.Sprintf("cart:%d", in.UserID)
	kv, _ := rdb.HGetAll(ctx, key).Result()
	cartInfos := make([]*__.CartInfo, 0, len(kv))
	if len(kv) > 0 {
		for pidStr, qtyStr := range kv {
			pid, _ := strconv.Atoi(pidStr)
			qty, _ := strconv.Atoi(qtyStr)
			cartInfos = append(cartInfos, &__.CartInfo{
				ProductID: int32(pid),
				UserID:    int32(in.UserID),
				Quantity:  int32(qty),
			})
		}
	} else {
		listdao := &dao.CartDao{}
		list, err := listdao.List(uint(in.UserID))
		if err != nil {
			return nil, fmt.Errorf("数据库查询失败" + err.Error())
		}
		for _, item := range list {
			_, _ = rdb.HSet(ctx, key, strconv.Itoa(int(item.ProductID)), item.Quantity).Result()
			cartInfos = append(cartInfos, &__.CartInfo{
				ProductID: int32(item.ProductID),
				UserID:    int32(item.UserID),
				Quantity:  int32(item.Quantity),
			})
		}
	}
	return &__.CartListResp{
		Code: 200,
		Msg:  "查询成功",
		Data: cartInfos,
	}, nil
}

func (s *Server) CartDelete(_ context.Context, in *__.CartDeleteReq) (*__.CartDeleteResp, error) {
	var ctx = context.Background()
	rdb := config.Rdb
	key := fmt.Sprintf("cart:%d", in.UserID)
	field := strconv.Itoa(int(in.ProductID))
	_, _ = rdb.HDel(ctx, key, field).Result()
	deletedao := &dao.CartDao{}
	err := deletedao.Remove(uint(in.UserID), uint(in.ProductID))
	if err != nil {
		return nil, fmt.Errorf("删除购物车商品失败: %v", err)
	}
	return &__.CartDeleteResp{
		Code: 200,
		Msg:  "删除成功",
	}, nil
}

func (s *Server) DelCart(_ context.Context, in *__.DelCartReq) (*__.DelCartResp, error) {
	var ctx = context.Background()
	rdb := config.Rdb
	key := fmt.Sprintf("cart:%d", in.UserID)
	_, _ = rdb.Del(ctx, key).Result()
	deldao := &dao.CartDao{}
	err := deldao.Clear(uint(in.UserID))
	if err != nil {
		return nil, fmt.Errorf("清空购物车失败: %v", err)
	}
	return &__.DelCartResp{
		Code: 200,
		Msg:  "清空成功",
	}, nil
}

func (s *Server) AuthorAuthCreate(_ context.Context, in *__.AuthorAuthCreateReq) (*__.AuthorAuthCreateResp, error) {
	authdao := &dao.AuthorAuthDao{}
	auth := &model.AuthorAuth{
		UserID:      uint(in.UserID),
		RealName:    in.RealName,
		AuthQualify: in.AuthQualify,
	}
	if err := authdao.Create(auth); err != nil {
		return nil, fmt.Errorf("入驻审核创建失败:" + err.Error())
	}
	return &__.AuthorAuthCreateResp{
		Code:   200,
		Msg:    "提交成功",
		AuthID: int32(auth.ID),
	}, nil
}

func (s *Server) AuthorAuthGet(_ context.Context, in *__.AuthorAuthGetReq) (*__.AuthorAuthGetResp, error) {
	authdao := &dao.AuthorAuthDao{}
	auth, err := authdao.Get(uint(in.UserID))
	if err != nil {
		return nil, fmt.Errorf("查询失败:" + err.Error())
	}
	info := &__.AuthorAuthInfo{
		Id:           int32(auth.ID),
		UserID:       int32(auth.UserID),
		RealName:     auth.RealName,
		AuthQualify:  auth.AuthQualify,
		Status:       int32(auth.Status),
		RejectReason: auth.RejectReason,
		Auditor:      auth.Auditor,
	}
	return &__.AuthorAuthGetResp{
		Code: 200,
		Msg:  "查询成功",
		Data: info,
	}, nil
}

func (s *Server) AuthorAuthUpdate(_ context.Context, in *__.AuthorAuthUpdateReq) (*__.AuthorAuthUpdateResp, error) {
	authdao := &dao.AuthorAuthDao{}
	if err := authdao.Update(uint(in.AuthID), int(in.Status), in.RejectReason, in.Auditor); err != nil {
		return nil, fmt.Errorf("更新失败:" + err.Error())
	}
	return &__.AuthorAuthUpdateResp{
		Code: 200,
		Msg:  "更新成功",
	}, nil
}
