package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/config"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/consumer"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/handlers"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/routes"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/services/mailer"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/services/mailer/templates"
	"github.com/Aditya-Nagpal/Cloud-File-Storage-System/services/notification-service/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables from .env file
	config.LoadConfig()

	// load templates once
	if err := templates.LoadTemplates(); err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}

	// context + graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create SQS client
	sqsClient, err := utils.NewSQSClient(ctx)
	if err != nil {
		log.Fatalf("failed to create sqs client: %v", err)
	}

	// mailer
	mailer := mailer.NewMailjetMailer()

	// processor
	proc := handlers.NewProcessor(mailer)

	// start consumer in background
	go consumer.Start(ctx, sqsClient, proc)

	// start http server (health)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	routes.SetupRoutes(r)

	srv := &httpServer{
		Engine: r,
		Port:   config.AppConfig.Port,
	}

	// run server in goroutine
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
	cancel()

	// allow some time for shutdown
	time.Sleep(2 * time.Second)
}

type httpServer struct {
	*gin.Engine
	Port string
}

func (s *httpServer) Run() error {
	return s.Engine.Run(s.Port) // s.Port already starts with :
}
