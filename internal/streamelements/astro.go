package streamelements

import (
	"fmt"
	"log"
	"net/url"

	"github.com/DaniruKun/tipfax/internal/config"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/securityguy/escpos"
)

const (
	TipsTopic           = "channel.tips"
	TipsModerationTopic = "channel.tips.moderation"
)

type Message struct {
	Type  string `json:"type"`
	Topic string `json:"topic"`
	Nonce string `json:"nonce"`
	Data  any    `json:"data"`
}

type Astro struct {
	cfg     *config.Config
	conn    *websocket.Conn
	printer *escpos.Escpos
}

func NewAstro(cfg *config.Config, printer *escpos.Escpos) *Astro {
	return &Astro{cfg: cfg, printer: printer}
}

func (a *Astro) Connect() error {
	u := url.URL{Scheme: "wss", Host: "astro.streamelements.com", Path: "/"}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	log.Println("Connected to Astro")

	a.conn = conn

	return nil
}

func (a *Astro) SubscribeTips() error {
	subscribeMessage := map[string]any{
		"type":  "subscribe",
		"nonce": uuid.New().String(),
		"data": map[string]any{
			"topic":      TipsTopic,
			"token":      a.cfg.SeJWTToken,
			"token_type": "jwt",
		},
	}

	if err := a.conn.WriteJSON(subscribeMessage); err != nil {
		log.Println("Error subscribing:", err)
		return err
	}

	log.Println("Subscribed to Astro topic:", TipsTopic)

	return nil
}

func (a *Astro) Listen() error {
	for {
		var msg Message
		err := a.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Error reading message:", err)
			return err
		}

		log.Printf("Received message: %+v", msg)

		// Handle different message types
		switch msg.Type {
		case "response":
			log.Println("Received response:", msg)
		case "message":
			log.Println("Received notification:", msg)
			// Process the notification based on the topic
			if msg.Topic == TipsTopic {
				a.handleTipMessage(msg)
			}
		}
	}
}

func (a *Astro) handleTipMessage(msg Message) {
	log.Println("🎉 NEW TIP RECEIVED! 🎉")

	// Parse the tip data
	data, ok := msg.Data.(map[string]any)
	if !ok {
		log.Println("Error parsing tip data")
		return
	}

	if donation, ok := data["donation"].(map[string]any); ok {
		username := "Unknown"
		amount := "0"
		currency := "USD"
		message := ""

		if user, ok := donation["user"].(map[string]any); ok {
			if name, ok := user["username"].(string); ok {
				username = name
			}
		}

		if amt, ok := donation["amount"].(float64); ok {
			amount = fmt.Sprintf("%.2f", amt)
		}
		if curr, ok := donation["currency"].(string); ok {
			currency = curr
		}

		if msg, ok := donation["message"].(string); ok {
			message = msg
		}

		status := "unknown"
		provider := "unknown"
		if statusVal, ok := data["status"].(string); ok {
			status = statusVal
		}
		if providerVal, ok := data["provider"].(string); ok {
			provider = providerVal
		}

		log.Printf("💰 Tip from %s: %s %s (via %s)", username, amount, currency, provider)
		log.Printf("📊 Status: %s", status)
		if message != "" {
			log.Printf("💬 Message: %s", message)
		}

		// Print to thermal printer if available
		if a.printer != nil {
			a.printer.Write(fmt.Sprintf("Tip from %s: %s %s", username, amount, currency))
			a.printer.LineFeed()
			a.printer.Write(fmt.Sprintf("Status: %s", status))
			a.printer.LineFeed()
			if message != "" {
				a.printer.Write(fmt.Sprintf("Message: %s", message))
				a.printer.LineFeed()
			}
			a.printer.PrintAndCut()
		}
	} else {
		log.Println("Error: Could not find donation data in tip message")
		log.Printf("Raw data: %+v", data)
	}
}

func (a *Astro) UnsubscribeTips() error {
	unsubscribeMessage := map[string]any{
		"type":  "unsubscribe",
		"nonce": uuid.New().String(),
		"data": map[string]any{
			"topic":      TipsTopic,
			"token":      a.cfg.SeJWTToken,
			"token_type": "jwt",
		},
	}

	if err := a.conn.WriteJSON(unsubscribeMessage); err != nil {
		log.Println("Error unsubscribing:", err)
	}

	log.Println("Unsubscribed from Astro topic:", TipsTopic)

	return nil
}

func (a *Astro) Disconnect() error {
	log.Println("Disconnecting from Astro")
	return a.conn.Close()
}
