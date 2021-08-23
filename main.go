package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"thinhtd4/config"
	"thinhtd4/controller/router"
	"thinhtd4/customlog"

	"thinhtd4/model"
	"thinhtd4/service/mongoservice"
	"thinhtd4/service/redisservice"

	"golang.org/x/net/http2"
)

func loadQuestionFromFileToDatabase() {
	file, er := ioutil.ReadFile("question.json")
	if er != nil {
		panic(er)
	}
	var qus []model.QuestionWithAnswer
	er = json.Unmarshal(file, &qus)
	if er != nil {
		panic(er)
	}

	for i, _ := range qus {
		er = mongoservice.MongoDB().InsertQuestion(context.Background(), &qus[i])
		if er != nil {
			panic(er)
		}
	}
}



/////////////  attention : add function call loadQuestionFromFileToDatabase() for just first time 
// 				to load question from file to db 
func main() {
	customlog.Info("Start main")
	customlog.Init()
	config.Init()
	mongoservice.Init()
	redisservice.Init()
	router.Init()

	// loadQuestionFromFileToDatabase()      uncomment this line when first run

	var httpServer = http.Server{
		Addr:    ":" + config.ServerConfig().Port,
		Handler: router.Router(),
	}

	var http2Server = http2.Server{}
	_ = http2.ConfigureServer(&httpServer, &http2Server)

	go func() {
		if err := httpServer.ListenAndServeTLS("./server.crt", "./server.key"); err != nil && errors.Is(err, http.ErrServerClosed) {
			customlog.Info("Start serving")
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	customlog.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		customlog.Info("Server forced to shutdown:")
		log.Fatal(err)
	}
	customlog.Info("Server exiting")
}
