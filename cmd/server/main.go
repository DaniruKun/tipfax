package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DaniruKun/tipfax/internal/config"
	"github.com/DaniruKun/tipfax/internal/fax"
	"github.com/DaniruKun/tipfax/internal/streamelements"
	"github.com/DaniruKun/tipfax/internal/web"
	"github.com/securityguy/escpos"
)

func main() {
	fmt.Println("Starting TipFax Server...")

	cfg := config.New()

	// Test printer connection
	printer, err := fax.NewPrinter(cfg.DevicePath)
	if err != nil {
		log.Printf("Warning: Failed to create printer: %v", err)
		log.Println("Continuing without printer...")
	} else {
		printer.SetConfig(escpos.ConfigEpsonTMT20II)
		printer.Write("TipFax Server Started!")
		printer.LineFeed()
		printer.PrintAndCut()
		log.Println("Printer test successful")
	}

	// Connect to StreamElements Astro
	astro := streamelements.NewAstro(cfg, printer)
	if err := astro.Connect(); err != nil {
		log.Fatalf("Failed to connect to StreamElements Astro: %v", err)
	}
	defer astro.Disconnect()

	if err := astro.SubscribeTips(); err != nil {
		log.Fatalf("Failed to subscribe to tips: %v", err)
	}
	defer astro.UnsubscribeTips()

	// Start HTTP server
	http.HandleFunc("/", web.StatusHandler(cfg, cfg.DevicePath))
	go func() {
		log.Printf("Starting HTTP server on http://localhost%s", cfg.ServerPort)
		if err := http.ListenAndServe(cfg.ServerPort, nil); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Start listening for messages in a goroutine
	go func() {
		log.Println("Starting to listen for tip messages...")
		if err := astro.Listen(); err != nil {
			log.Printf("Error listening for messages: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("TipFax Server is running. Press Ctrl+C to stop.")
	log.Printf("Web interface available at: http://localhost%s", cfg.ServerPort)

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down TipFax Server...")
}
