package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-Sakura/SMS/model/mysql"
)

func Rand(lenth int) string {
	var code string
	rand.Seed(time.Now().Unix())
	for i := 0; i < lenth; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

func Send(db *sql.DB, mobile, code string) error {
	id, err := mysql.GetSMServiceId(db)
	if err != nil {
		return err
	}

	smservice, err := mysql.GetSMService(db, id)
	if err != nil {
		return err
	}

	switch smservice.Provider {
	case "aliyun":
		client := NewAliClient(smservice.AccessKeyId, smservice.AccessSecret, smservice.TemplateCode, smservice.SignName)
		err := client.Send(mobile, code)
		if err != nil {
			return err
		}
	case "tianyan":
		client := NewTianYanClient(smservice.AppCode, smservice.TemplateCode)
		err := client.Send(mobile, code)
		if err != nil {
			return err
		}
	}

	b, _ := strconv.Atoi(smservice.Balance)
	b--
	balance := strconv.Itoa(b)
	err = mysql.UpDateBalance(db, balance, id)
	if err != nil {
		return errors.New("errUpDate")
	}

	return nil
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")

			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")

			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
