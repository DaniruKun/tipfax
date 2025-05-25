package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DaniruKun/tipfax/internal/config"
	"github.com/DaniruKun/tipfax/internal/fax"
	"github.com/DaniruKun/tipfax/internal/streamelements"
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

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down TipFax Server...")
}
