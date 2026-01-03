package pkg

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
)

func NewAlipayClient() (*alipay.Client, error) {
	privateKey := "MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCbTU0brJeQvWZMiBK8gusVWV7gkUFHm49Cc3h3sOZZqClgQ6YnxqoIFi4LkdnL4XI5tHVkwF7zbW10Vo+bQ1Vy2cXoXysq0+ERTk90+5eZe9WgOgQw7nEl5bpsI1gcVQlMc4jHzLdzhqtxV7V+06xhB8E7hx7aKcXK6qPWzgPjD1vw3Jh7Kh8sc3xwWurvuzBv/zDd6zPLQL3mrred7tSDWaYxHlygsqPuasLOp8s7OC1eHPm4Q8hrhhLlXJ93bcfK9rRf+11ovsM8af7ViOEkBZ8sSoT+zd7eAswXudPImxK+VJhfE+HbkPmEBcNlamhOa+mP7z5dbTYn4ORTKg3zAgMBAAECggEAJlY0vKokWBSJFkdY6LQguApxw1nYFYcvHCZJDLLcijFa1Wqdr5/5BToEb9K3Qv8KJXiIqjBawBi5NkjP9nHuvAVKN4yWqWHYY19DegtJZxgMqhroQfP6dnJ5TIyFCalsPDUhBMFiX+BUngwav44dNW6mor0+PnilXOwDOdltWDVl7yGj+hsrmAr4B11BrONtJQDm48lDU5Exl28lFZBuz7q9vlo431IF85ExjkahIS1XwAhc6jgFCbsZFTe0TUwcTWKaY6HGe0Z0RHuySoLwcvmmhajBlJpEj7gwIkqgUjBiq9lx/Kvw8mXKAbEWNeX7Ux9wG+gJh9+TpymCkZYyMQKBgQDNMrnha67xNtZkygcWSAipdYH7f+ogfaP3R2Z7zu66JcvGdKHW1laRbwnLkj+QkwXpALjw7jeAQGxtGS5AsDpHTzXJb1L7L51JsnUGvSLecVDKGuuwOXqHDMKDMMpGKTSCV3ShfpJgzUY73atyf+tqXhXQvtCZfGdx90CpsHZY1wKBgQDBwDMPPE7gj6KZFZzjp2Ulx+o8uvOjmsMFFdLX4ojpNrKYCVP+cIdYsvp417NWy4ojfhTN+mxACe5XJcOQ7Jo8w80zYODvzKB8Lrdu7hDUrvBOgweO5wBYAenrgGJs6t6ngElKjV4w6v4tmMIvnnGFsN7Dyg9V2wYG7Hgy3U1ERQKBgQCx+7bJNQr4BNWvdKDzDaYdvmPxTsE4T2JCYTceWp4s7g2zZITgAOfzm7mGTqM99pigwtSnfii74KVrd/TrfS//nFAOGbaDU4h9XQIuxy0Qfn9R3knif+isbT/mZRJ+Cs2V5N+wGEZFGqg50wsb9KKwj00i0+/BwetEKe93gC0W3wKBgGE8z6NW1hNXovgHY8zRRy817PX7sakrU7Lqp/2XALVisTEihVgOK681bAVX4/asgjCb518Uzl05Xre4CTVjuWjDr+mYNmvDG8wXOhJfQm0rOwl8Mz/h6UdB9p8tuLgHDCWueZoD5wDP/y7tGpABieHZyYMjlpy1Joo1BYIplMytAoGBAJad7PWmpzzk6RioUDmK3BtMiQBHlGyqNI/hllDu+fyZgZxEnAXMAnQjpnbnyLZj7GQdgPZNV5NYd3e84qBCxQiOY69AhS0ltMJJML5ZPqIXnNtKrjwnzfJlzvWOxS+SwV6/vp5602MBjTDgAiXPNRW55dblMZdcGshFN8CdeVq3"
	appId := "9021000151640003"
	return alipay.New(appId, privateKey, false)
}
func Alipay(Subject, OutTradeNo, TotalAmount string) (string, error) {
	client, err := NewAlipayClient()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var p = alipay.TradeWapPay{}
	p.NotifyURL = "https://3da358.r19.cpolar.top/alipay/notify"
	p.ReturnURL = "https://3da358.r19.cpolar.top/return"
	p.Subject = Subject         //"标题"
	p.OutTradeNo = OutTradeNo   //"传递一个唯一单号"
	p.TotalAmount = TotalAmount //"10.00"
	p.ProductCode = "QUICK_WAP_WAY"

	url, err := client.TradeWapPay(p)
	if err != nil {
		fmt.Println(err)
	}

	// 这个 payURL 即是用于打开支付宝支付页面的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()
	fmt.Println(payURL)
	return payURL, nil
}
