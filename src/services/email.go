package services

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/myroslavve/genesis-test-case/src/api"
	"github.com/myroslavve/genesis-test-case/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func InitEmailService() {
	collection = db.GetCollection("subscriptions")
	log.Println("Email service initialized")
}

// sanitizeEmail removes any CR or LF characters from the email address
func sanitizeEmail(email string) string {
	return strings.ReplaceAll(strings.ReplaceAll(email, "\n", ""), "\r", "")
}

func SendEmail(to string, rate float64) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Sanitize email addresses
	from = sanitizeEmail(from)
	to = sanitizeEmail(to)

	log.Printf("Sending email from %s to %s via %s:%s\n", from, to, smtpHost, smtpPort)

	subject := "Subject: Current USD to UAH Exchange Rate"
	body := fmt.Sprintf("The current USD to UAH exchange rate is %f.", rate)

	var msg strings.Builder
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString(subject + "\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg.String()))
	if err != nil {
		log.Printf("Failed to send email to %s: %v\n", to, err)
		return err
	}
	log.Printf("Email sent to %s\n", to)
	return nil
}

func SendExchangeRates(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		rate, err := api.FetchExchangeRate()
		if err != nil {
			log.Println("Failed to fetch exchange rate:", err)
			continue
		}
		log.Printf("Fetched exchange rate: %f\n", rate)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			log.Println("Failed to fetch subscriptions:", err)
			cancel()
			continue
		}

		for cursor.Next(ctx) {
			var result bson.M
			if err := cursor.Decode(&result); err != nil {
				log.Println("Failed to decode subscription:", err)
				continue
			}

			email, ok := result["email"].(string)
			if !ok {
				log.Println("Invalid email format in database")
				continue
			}

			if err := SendEmail(email, rate); err != nil {
				log.Println("Failed to send email to", email, ":", err)
			} else {
				log.Println("Sent exchange rate to", email)
			}
		}

		cursor.Close(ctx)
		cancel()
	}
}
