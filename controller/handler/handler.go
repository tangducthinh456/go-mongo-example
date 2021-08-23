package handler

import (
	"errors"
	"fmt"
	// "io/ioutil"

	// "fmt"

	// "fmt"
	"net/http"
	"net/url"

	// "os"
	// "thinhtd4/config"
	"thinhtd4/config"
	"thinhtd4/customlog"
	"thinhtd4/model"
	"thinhtd4/service/mongoservice"
	"thinhtd4/service/redisservice"

	"github.com/gin-gonic/gin"
)

func HandlleLoginGET(c *gin.Context) {
	customlog.Info("Get login endpoint")
	c.HTML(http.StatusOK, "login-form.html", gin.H{"title": "login"})
}

func HandleRegisterGET(c *gin.Context) {
	customlog.Info("Get register endpoint")
	c.HTML(http.StatusOK, "register-form.html", gin.H{"title": "register"})
}

func HandleLoginPOST(c *gin.Context) {
	customlog.Info("Post login endpoint")
	var user *model.UserDto
	err := c.Bind(&user)
	if err != nil {
		customlog.Err(errors.New("Post register endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}
	userInfo, err := mongoservice.MongoDB().FindUser(c, user.Username)
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
	if !redisservice.CheckPasswordHash(user.Password, userInfo.EncryptPassword) {
		customlog.Info("Post login endpoint password not match")
		c.String(http.StatusUnauthorized, "Your username or password is not corrent. Try again")
		return
	}

	token := model.CookieResp{
		UserID:  userInfo.UserID.Hex(),
		TimeOut: 3600,
	}

	err = redisservice.SaveTokenToRedis(token)
	if err != nil {
		customlog.Err(errors.New("Post register endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}

	c.SetCookie("user_id", token.UserID, int(token.TimeOut), "/", config.ServerConfig().Host, true, true)
	c.SetCookie("username", userInfo.Username, int(token.TimeOut), "/", config.ServerConfig().Host, true, true)

	location := url.URL{Path: "/survey"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func HandleRegisterPOST(c *gin.Context) {
	customlog.Info("Post register endpoint")
	var user *model.UserDto
	err := c.Bind(&user)
	if err != nil {
		customlog.Err(errors.New("Post register endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}

	// validate if exist username
	userInfoFromDB, err := mongoservice.MongoDB().FindUser(c, user.Username)
	fmt.Println(user.Username)
	fmt.Println(userInfoFromDB.Username)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			customlog.Info("Post register endpoint : " + err.Error())
			c.String(http.StatusUnauthorized, "Your username or password is not corrent. Try again")
			return
		}
	} else {
		customlog.Info("Post register endpoint username already exist")
		c.String(http.StatusUnauthorized, "Your username already exist. Try new one")
		return
	}

	pwdEncrypt, err := redisservice.HashPassword(user.Password)
	if err != nil {
		customlog.Err(errors.New("Post register endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}

	userInfo := &model.UserInfo{
		Username:        user.Username,
		EncryptPassword: pwdEncrypt,
	}

	err = mongoservice.MongoDB().InsertUser(c, userInfo)
	if err != nil {
		customlog.Err(errors.New("Post register endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}
	c.HTML(http.StatusCreated, "after-register.html", gin.H{"title": "Login"})
}

func HandleGetListQuestion(c *gin.Context) {
	customlog.Info("Get question endpoint")
	list, er := mongoservice.MongoDB().GetQuestionsList(c)
	if er != nil {
		customlog.Err(errors.New("handle get list error : " + er.Error()))
		c.String(http.StatusInternalServerError, "Please try again later")
		return
	}
	var questionDto []model.QuestionDto
	for _, v := range list {
		questionDto = append(questionDto, model.QuestionDto{
			QuestionID: v.QuestionID,
			Question:   v.Question,
		})
	}
	c.JSON(http.StatusOK, questionDto)
}

func HandleUpdateUserQuestionDone(c *gin.Context) {
	customlog.Info("Post login endpoint")
	// answer := model.AnswerUpdate{}
	// answer.Answer = make([]int, 0)
	var answer *model.AnswerUpdate = &model.AnswerUpdate{}
	answer.Answer = make([]int, 0)
	err := c.BindJSON(&answer)
	if err != nil {
		customlog.Err(errors.New("PUT answer endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}
	// bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	// fmt.Println(string(bodyBytes))
	fmt.Println(answer)
	cookie, _ := c.Request.Cookie("user_id")
	userID := cookie.Value

	err = mongoservice.MongoDB().UpdateUserQuestionDone(c, userID, answer.QuestionID, answer.Answer)
	if err != nil {
		customlog.Err(errors.New("PUT answer endpoint error : " + err.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}
	c.String(http.StatusAccepted, "Ok")
}

func HandleSubmit(c *gin.Context) {
	customlog.Info("Submit survey")
	cookie, _ := c.Request.Cookie("user_id")
	userID := cookie.Value
	er := mongoservice.MongoDB().UpdateUserHaveDoneSurvey(c, userID, true)
	if er != nil {
		customlog.Err(errors.New("PUT submit all endpoint error : " + er.Error()))
		c.String(http.StatusInternalServerError, "Error, please try again later")
		return
	}
	list, er := mongoservice.MongoDB().GetQuestionsList(c)
	if er != nil {
		customlog.Err(errors.New("handle get list error : " + er.Error()))
		c.String(http.StatusInternalServerError, "Please try again later")
		return
	}
	c.JSON(http.StatusOK, list)
}

func HandleSurvey(c *gin.Context) {
	customlog.Info("Access to survey")
	c.HTML(200, "survey.html", gin.H{})
}
