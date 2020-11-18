package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var log = logrus.New()

func main() {
	r := gin.Default()

	viper.SetConfigName("backend")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if viper.Get("mode").(string) == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DisableConsoleColor()
	//f, _ := os.OpenFile("../log/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	//gin.DefaultWriter = io.MultiWriter(f)
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("../log/gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		PrettyPrint:      true,
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.Get("port").(int)),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server started.")

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
