package middleware

import (
	"errors"
	"net/http"
	"net/url"

	"thinhtd4/customlog"
	"thinhtd4/service/mongoservice"
	"thinhtd4/service/redisservice"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		location := url.URL{Path: "/login"}
		cookie, _ := c.Request.Cookie("user_id")

		if cookie != nil {
			check, err := redisservice.IsTokenInRedis(cookie.Value)
			if err != nil {
				customlog.Err(errors.New("Auth middleware error : " + err.Error()))
				c.Redirect(http.StatusFound, location.RequestURI())
			}
			if !check {
				customlog.Info("Invalid access token : ")
				c.Redirect(http.StatusFound, location.RequestURI())
			}

			cookie, _ = c.Request.Cookie("username")
			userInfo, err := mongoservice.MongoDB().FindUser(c, cookie.Value)
			if err != nil {
				if err.Error() == "mongo: no documents in result" {
					customlog.Info("Post login endpoint : " + err.Error())
					c.String(http.StatusUnauthorized, "Your username or password is not corrent. Try again")
					return
				}
				customlog.Err(errors.New("Post register endpoint error : " + err.Error()))
				c.String(http.StatusInternalServerError, "Error, please try again later")
				return
			}
			if userInfo.HaveDoneSurvey {
				c.String(http.StatusOK, "You have done your survey, no more test. You are free now.")
				c.Abort()
				return
			}

		} else {
			c.Redirect(http.StatusFound, location.RequestURI())
		}

	}
}
