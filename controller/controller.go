package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-Sakura/SMS/model/mysql"
	"github.com/my-Sakura/SMS/utils"
)

var (
	errSend   = errors.New("Send failed")
	errUpDate = errors.New("errUpDate")
)

type SMSController struct {
	DB     *sql.DB
	Length int
}

func NewSMSController(db *sql.DB) *SMSController {
	return &SMSController{
		DB:     db,
		Length: 6,
	}
}

func (s *SMSController) RegistRouter(r gin.IRouter) {
	err := mysql.CreateDatabase(s.DB)
	if err != nil {
		log.Println(err)
	}

	mysql.CreateTableMessage(s.DB)
	mysql.CreateTableSMService(s.DB)

	r.Static("/", "./view")
	//API for suppliers
	r.POST("/add", s.Add)
	r.POST("/delete", s.Delete)
	r.POST("/update", s.UpDate)
	r.POST("/get", s.Get)
	//API for customers
	r.POST("/send", s.Send)
	r.POST("/check", s.Check)
}

func (s *SMSController) Add(c *gin.Context) {
	var (
		req struct {
			Provider     string `json:"provider"`
			AppCode      string `json:"appcode"`
			Balance      string `json:"balance"`
			TemplateCode string `json:"templatecode"`
			SignName     string `json:"singname"`
			AccessKeyId  string `json:"accesskeyid"`
			AccessSecret string `json:"accesssecret"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	switch req.Provider {
	case "aliyun":
		err := mysql.InsertALiYunSMService(s.DB, req.Provider, req.AccessKeyId, req.AccessSecret, req.TemplateCode, req.SignName, req.Balance)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
			return
		}
	case "tianyan":
		err := mysql.InsertTianYanSMService(s.DB, req.Provider, req.AppCode, req.TemplateCode, req.Balance)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (s *SMSController) Delete(c *gin.Context) {
	var (
		req struct {
			Id string `json:"id"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteSMService(s.DB, req.Id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (s *SMSController) UpDate(c *gin.Context) {
	var (
		req struct {
			Id      string `json:"id"`
			Balance string `json:"balance"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.UpDateBalance(s.DB, req.Balance, req.Id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": errUpDate})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (s *SMSController) Get(c *gin.Context) {
	var (
		req struct {
			Id string `json:"id"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	smservice, err := mysql.GetSMService(s.DB, req.Id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"smservice": smservice,
	})
}

func (s *SMSController) Send(c *gin.Context) {
	//The front end needs to send a unique Id
	var (
		req struct {
			Mobile string `json:"mobile"`
			Id     string `json:"id"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	code := utils.Rand(s.Length)
	err = mysql.InsertMessage(s.DB, req.Id, req.Mobile, code)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errSend})
		return
	}

	err = utils.Send(s.DB, req.Mobile, code)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errSend})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func (s *SMSController) Check(c *gin.Context) {
	var (
		req struct {
			Code string `json:"code"`
			Id   string `json:"id"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	msg, err := mysql.GetMessage(s.DB, req.Id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if msg.Code != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{"err": "验证码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "验证成功"})

	err = mysql.DeleteMessage(s.DB, req.Id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
}
