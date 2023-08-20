package ginx

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	schema "github.com/easyai-io/easyai-platform/pkg/apischema"
	"github.com/easyai-io/easyai-platform/pkg/contextx"
	"github.com/easyai-io/easyai-platform/pkg/errors"
	"github.com/easyai-io/easyai-platform/pkg/logger"
	"github.com/easyai-io/easyai-platform/pkg/util/json"
)

const (
	prefix = "easyai-platform"
	// ReqBodyKey ...
	ReqBodyKey = prefix + "/req-body"
	// ResBodyKey ...
	ResBodyKey = prefix + "/res-body"
)

// GetToken Get jwt token from header (X-Authorization-Token: Bearer xxx)
func GetToken(c *gin.Context) string {
	token := c.GetHeader("X-Authorization-Token")
	prefix := "Bearer "
	if token != "" && strings.HasPrefix(token, prefix) {
		token = token[len(prefix):]
	}
	if token == "" {
		token, _ = c.Cookie("easyai_token")
	}
	return token
}

// GetBodyData Get body data from context
func GetBodyData(c *gin.Context) []byte {
	if v, ok := c.Get(ReqBodyKey); ok {
		if b, ok := v.([]byte); ok {
			return b
		}
	}
	return nil
}

// ParseParamID Param returns the value of the URL param
func ParseParamID(c *gin.Context, key string) uint32 {
	val := c.Param(key)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0
	}
	return uint32(id)
}

// ParseJSON Parse body json data to struct
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request json failed: %s", err.Error()))
	}
	return nil
}

// ParseQuery Parse query parameter to struct
// for gin/binding.Query, the obj' field should have a `form` tag
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request query failed: %s", err.Error()))
	}
	return nil
}

// QueryInt get int query arg
func QueryInt(c *gin.Context, key string) int {
	val, ok := c.GetQuery(key)
	if !ok {
		return -1
	}
	if v, err := strconv.Atoi(val); err == nil {
		return v
	}
	return -1
}

// QueryString get string query arg
func QueryString(c *gin.Context, key string) string {
	return c.Query(key)
}

// QueryStringSplit get string slice query arg
func QueryStringSplit(c *gin.Context, key, sep string) []string {
	v, ok := c.GetQuery(key)
	if !ok {
		return nil
	}
	return strings.Split(v, sep)
}

// QueryBool get bool query arg
func QueryBool(c *gin.Context, key string) bool {
	val, ok := c.GetQuery(key)
	if !ok {
		return false
	}
	if v, err := strconv.ParseBool(val); err == nil {
		return v
	}
	return false
}

// ParseForm Parse body form data to struct
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.Wrap400Response(err, fmt.Sprintf("Parse request form failed: %s", err.Error()))
	}
	return nil
}

// ResOK success with message ok
func ResOK(c *gin.Context) {
	ResJSONBody(c, http.StatusOK, 200, string(schema.OKStatus), nil)
}

// ResSuccess success with data
func ResSuccess(c *gin.Context, v interface{}) {
	ResJSONBody(c, http.StatusOK, 200, "", v)
}

// ResList data with list object
func ResList(c *gin.Context, v interface{}) {
	ResJSONBody(c, http.StatusOK, 200, "", v)
}

// ResPage response pagination data object
func ResPage(c *gin.Context, v interface{}, pr *schema.PaginationResult) {
	list := schema.ListResult{
		List:       v,
		Pagination: pr,
	}
	ResJSONBody(c, http.StatusOK, 200, "", list)
}

// ResJSONBody response json body
func ResJSONBody(c *gin.Context, httpCode, bizCode int, msg string, data interface{}) {
	v := schema.RespData{
		Code:    bizCode,
		Message: msg,
		Data:    data,
	}

	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	c.Set(ResBodyKey, buf)
	c.Data(httpCode, "application/json; charset=utf-8", buf)
	c.Abort()
}

// ResError Response error object and parse error status code
func ResError(c *gin.Context, err error, status ...int) {
	ctx := c.Request.Context()
	var res *errors.ResponseError

	if err != nil {
		if e, ok := err.(*errors.ResponseError); ok {
			res = e
		} else {
			res = errors.UnWrapResponse(errors.ErrInternalServer)
			res.ERR = err
		}
	} else {
		res = errors.UnWrapResponse(errors.ErrInternalServer)
	}

	if res.Code <= 0 {
		res.Code = 400
	}

	if len(status) > 0 {
		res.Status = status[0]
	}

	if err := res.ERR; err != nil {
		if res.Message == "" {
			res.Message = err.Error()
		}

		if status := res.Status; status >= 400 && status < 500 {
			logger.WithContext(ctx).Warnf(err.Error())
		} else if status >= 500 {
			logger.WithContext(logger.NewStackContext(ctx, err)).Errorf(err.Error())
		}
	}

	eitem := schema.ErrorItem{
		Code:    res.Code,
		Message: res.Message,
	}
	ResJSONBody(c, res.Status, res.Code, res.Message, eitem)
}

// ResOctetStream for binary response
func ResOctetStream(c *gin.Context, httpCode int, buf []byte) {
	// c.Set(ResBodyKey, buf)  // we do not log for this
	c.Data(httpCode, "application/octet-stream", buf)
	c.Abort()
}

// UserIDFromCtx get user id
func UserIDFromCtx(ctx context.Context) string {
	if userID := contextx.FromUserUID(ctx); userID != "" {
		return userID
	}
	if g, ok := ctx.(*gin.Context); ok {
		if userID := contextx.FromUserUID(g.Request.Context()); userID != "" {
			return userID
		}
	}
	return ""
}

// SetCookieIfNotExist set cookie if not exist
func SetCookieIfNotExist(c *gin.Context, name, value string, maxAge int, path string, secure, httpOnly, rootDomain bool) {
	if _, err := c.Cookie(name); err == nil {
		return
	}
	domain := c.Request.Host
	if rootDomain {
		if parts := strings.Split(domain, "."); len(parts) > 2 {
			// .开头表示可以被子域名访问
			domain = fmt.Sprintf(".%s.%s", parts[len(parts)-2], parts[len(parts)-1])
		}
	}
	c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}
