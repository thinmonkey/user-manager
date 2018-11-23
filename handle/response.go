package handle

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"digitalkey-manager/utils"
)

const (
	SUCCESS                     = 1000
	SYSTEM_ERROE                = 1500
	QUERY_PARAM_ERROR           = 1400
	RESOURCE_NOT_FOUNT          = 1405
)

type Response struct {
	CustomeError
	Data interface{} `json:"data,omitempty"` //omitempty有值就输出，没值则不输出
}

type AdminResponse struct {
	Error *CustomeError `json:"error,omitempty"`
	Data  interface{}   `json:"data,omitempty"` //omitempty有值就输出，没值则不输出
}

type CustomeErrorContainer map[int]string

var customeErrorContainer = CustomeErrorContainer{
	SUCCESS:                     "success",
	QUERY_PARAM_ERROR:           "请求参数有误",
	RESOURCE_NOT_FOUNT:          "访问资源不存在",
	SYSTEM_ERROE:                "服务内部错误",
}

type CustomeError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func (customeError CustomeError) Error() string {
	return customeError.Message
}

func NewCustomeError(code int) CustomeError {
	return CustomeError{Code: code, Message: customeErrorContainer[code]}
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	utils.CheckErr(err)
	var customeError CustomeError
	if err == nil {
		customeError = NewCustomeError(SUCCESS)
		c.JSON(http.StatusOK, Response{
			CustomeError: customeError,
			Data:         data,
		})
	} else {
		switch err.(type) {
		case CustomeError:
			customeError = err.(CustomeError)
			break
		default:
			customeError = NewCustomeError(SYSTEM_ERROE)
			break
		}
		c.JSON(http.StatusOK, Response{
			CustomeError: customeError,
		})
	}

}

