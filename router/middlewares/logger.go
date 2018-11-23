package middlewares

import (
	"bytes"

	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"
	"github.com/thinmonkey/user-manager/utils/log"
)

// 2016-09-27 09:38:21.541541811 +0200 CEST
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700]
// "GET /apache_pb.gif HTTP/1.0" 200 2326
// "http://www.example.com/start.html"
// "Mozilla/4.08 [en] (Win98; I ;Nav)"
const DefaultMemory = 32 * 1024 * 1024

var timeFormat = "2006-01-02 15:04:05.000"

func LoggerMiddlerware() gin.HandlerFunc {
	return func(context *gin.Context) {
		log := log.GetLogrus()
		url := context.Request.URL
		method := context.Request.Method
		path := url.Path
		host := context.Request.Host
		clientUserAgent := context.Request.UserAgent()
		clientIp := GetIP(context)

		var requestBody interface{}
		requestBody = GetRequestBody(context)
		requestHeaders := GetHeaders(context.Request.Header)

		start := time.Now()
		context.Next()
		stop := time.Since(start)
		requestTime := int(math.Ceil(float64(stop.Nanoseconds()) / 1000.0))
		statusCode := context.Writer.Status()

		dataLength := context.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logrus.NewEntry(log).WithFields(logrus.Fields{
			"hostname":    host,
			"statusCode":  statusCode,
			"requestTime": requestTime, // time to process
			"clientIP":    clientIp,
			"headers":     requestHeaders,
			"method":      method,
			"path":        path,
			"dataLength":  dataLength,
			"userAgent":   clientUserAgent,
			"body":        requestBody,
		})

		if len(context.Errors) > 0 {
			entry.Error(context.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("currentTime:[%s] success requestTime[%dms]", time.Now().Format(timeFormat), requestTime)
			if statusCode > 499 {
				entry.Error(msg)
			} else if statusCode > 399 {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}

// GetHeaders ...
func GetHeaders(head http.Header) map[string]string {
	hdr := make(map[string]string, len(head))
	for k, v := range head {
		hdr[k] = v[0]
	}
	return hdr
}

// GetIP ...
func GetIP(c *gin.Context) string {
	ip := c.ClientIP()
	return ip
}

// GetMultiPartFormValue ...
func GetMultiPartFormValue(c *http.Request) interface{} {
	var requestBody interface{}

	multipartForm := make(map[string]interface{})
	if err := c.ParseMultipartForm(DefaultMemory); err != nil {
		// handle error
	}
	if c.MultipartForm != nil {
		for key, values := range c.MultipartForm.Value {
			multipartForm[key] = strings.Join(values, "")
		}

		for key, file := range c.MultipartForm.File {
			for k, f := range file {
				formKey := fmt.Sprintf("%s%d", key, k)
				multipartForm[formKey] = map[string]interface{}{"filename": f.Filename, "size": f.Size}
			}
		}

		if len(multipartForm) > 0 {
			requestBody = multipartForm
		}
	}
	return requestBody
}

// GetFormBody ...
func GetFormBody(c *http.Request) interface{} {
	var requestBody interface{}

	form := make(map[string]string)
	if err := c.ParseForm(); err != nil {
		// handle error
	}
	for key, values := range c.PostForm {
		form[key] = strings.Join(values, "")
	}
	if len(form) > 0 {
		requestBody = form
	}

	return requestBody
}

// GetRequestBody ...
func GetRequestBody(c *gin.Context) interface{} {
	//multiPartFormValue := GetMultiPartFormValue(c.Request)
	//if multiPartFormValue != nil {
	//	return multiPartFormValue
	//}
	//
	//formBody := GetFormBody(c.Request)
	//if formBody != nil {
	//	return formBody
	//}

	method := c.Request.Method
	if method == "GET" {
		return nil
	}
	contentType := c.ContentType()
	body := c.Request.Body
	var model interface{}
	bodyContent, err := ioutil.ReadAll(body)
	if err != nil {
		return model
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyContent))

	switch contentType {
	case binding.MIMEJSON:
		json.Unmarshal(bodyContent, &model)
		return model
	default:
		model = string(bodyContent)
		return model
	}
}
