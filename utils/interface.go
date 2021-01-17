package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

type SMService interface {
	Send(mobile, code string) error
}

//ALiYun  Service Provider
type AliClient struct {
	Scheme       string
	SignName     string
	TemplateCode string
	Regioned     string
	AccessKeyId  string
	AccessSecret string
}

//TianYan Service Provider
type TianYanClient struct {
	TemplateCode string
	AppCode      string
	Method       string
	Url          string
}

func NewAliClient(accessKeyId, accessSecret, templateCode, signName string) *AliClient {
	return &AliClient{
		Scheme:       "https",
		Regioned:     "cn-hangzhou",
		SignName:     signName,
		TemplateCode: templateCode,
		AccessKeyId:  accessKeyId,
		AccessSecret: accessSecret,
	}
}

func NewTianYanClient(appCode, templateCode string) *TianYanClient {
	return &TianYanClient{
		TemplateCode: templateCode,
		AppCode:      appCode,
		Method:       "POST",
		Url:          "https://smssend.shumaidata.com/sms/send?receive=%s&tag=%s&templateId=%s",
	}
}

func (a *AliClient) Send(mobile, code string) error {
	client, err := dysmsapi.NewClientWithAccessKey(a.Regioned, a.AccessKeyId, a.AccessSecret)
	if err != nil {
		fmt.Print(err.Error())
	}

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = a.Scheme
	request.SignName = a.SignName
	request.TemplateCode = a.TemplateCode
	request.TemplateParam = "{\"code\":" + code + "}"
	request.PhoneNumbers = mobile

	response, err := client.SendSms(request)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("response is %#v\n", response)
	return nil
}

func (t *TianYanClient) Send(mobile, code string) error {
	client := &http.Client{}

	url := fmt.Sprintf(t.Url, mobile, code, t.TemplateCode)
	req, err := http.NewRequest(t.Method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "APPCODE"+" "+t.AppCode)
	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response is %v\n", string(b))

	return nil
}
