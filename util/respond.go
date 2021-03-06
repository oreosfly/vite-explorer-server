package util

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vitelabs/vite-explorer-server/vitelog"
)

type responseData interface {
	ToResponse () gin.H
}

func Respond (c *gin.Context, data responseData, msg string, err error, code int) {
	var resData gin.H
	if data != nil {
		resData = data.ToResponse()
	}
	var errStr = ""

	if err != nil {
		errStr = err.Error()
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.JSON(200, gin.H{
		"code": code, // success
		"msg": msg,
		"data": resData,
		"error": errStr,
	})


}

func RespondSuccess (c *gin.Context, data responseData, msg string)  {
	Respond(c, data, msg, nil, 0)
}

func RespondFailed (c *gin.Context, code int,err error , msg string)  {
	vitelog.Logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url": c.Request.URL,
		"code": code,
	}).Error(err)

	Respond(c, nil, msg, err, code)
}


func RespondError (c *gin.Context, statusCode int, err error )  {
	c.String(statusCode,  err.Error())
}
