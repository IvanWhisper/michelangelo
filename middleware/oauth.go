package middleware

import (
	"context"
	"fmt"
	mlog "github.com/IvanWhisper/michelangelo/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

func OAuth(authurl, clientid, secret string,getToken func(resp []byte) (interface{},error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// catch a bear token
		bearToken := c.Request.Header.Get("Authorization")
		var token string
		if vs := strings.Split(bearToken, " "); len(vs) > 1 {
			token = vs[1]
		}
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// check token with auth server
		url := fmt.Sprintf("%s/check?access_token=%s", authurl, token)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tk,err:=getToken(bytes)
		if err != nil {
			mlog.Error("token unknow"+err.Error())
		} else {
			newCtx:=context.WithValue(c.Request.Context(),"tokeninfo",tk)
			c.Request=c.Request.WithContext(newCtx)
		}
		c.Next()
	}
}
