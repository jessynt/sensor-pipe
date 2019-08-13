package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sensorsanalytics "github.com/sensorsdata/sa-sdk-go"
	"github.com/sensorsdata/sa-sdk-go/consumers"
	"github.com/sensorsdata/sa-sdk-go/structs"

	"sensor-pipe/handler"
)

type StdoutConsumer struct{}

func (s StdoutConsumer) Send(data structs.EventData) error {
	itemData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	log.Println(string(itemData))
	return nil
}

func (s StdoutConsumer) Flush() error {
	return nil
}

func (s StdoutConsumer) Close() error {
	return nil
}

func (s StdoutConsumer) ItemSend(item structs.Item) error {
	return nil
}

func main() {
	var err error

	_, debugEnabled := os.LookupEnv("DEBUG")
	saServerURL, exists := os.LookupEnv("SA_SERVER_URL")
	if !exists && !debugEnabled {
		panic("SA_SERVER_URL is required")
	}
	saProject, exists := os.LookupEnv("SA_PROJECT")
	if !exists && !debugEnabled {
		panic("SA_PROJECT required")
	}

	var consumer consumers.Consumer

	consumer = &StdoutConsumer{}
	if !debugEnabled {
		consumer, err = sensorsanalytics.InitBatchConsumer(saServerURL, 10, 1000*10)
		if err != nil {
			panic(err)
		}
	}

	sa := sensorsanalytics.InitSensorsAnalytics(consumer, saProject, false)

	var route *gin.Engine
	{
		gin.SetMode(gin.ReleaseMode)
		route = gin.New()
		route.Use(gin.Logger(), gin.Recovery())
		route.GET("/ping", handler.MakePingHandler())
		route.POST("/track", handler.MakeTrackHandler(sa))
		route.POST("/track-signup", handler.MakeTrackSignupHandler(sa))
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "0.0.0.0", 80),
		Handler: route,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGABRT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	sa.Close()

	log.Println("Server exiting")
}
