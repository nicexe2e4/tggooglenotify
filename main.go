package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	tb "gopkg.in/telebot.v3"
)

var telegramToken string
var chatID int64
var calendarID string
var cronExpression string
var userTags map[string]string

func init() {
	telegramToken = os.Getenv("TELEGRAM_TOKEN")
	calendarID = os.Getenv("CALENDAR_ID")
	cronExpression = os.Getenv("CRON_EXPRESSION")

	var err error
	chatIDStr := os.Getenv("CHAT_ID")
	if chatIDStr == "" {
		log.Fatalf("CHAT_ID environment variable is not set")
	}
	chatID, err = strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid chat ID: %v", err)
	}

	userTagsJSON := os.Getenv("USER_TAGS")
	if userTagsJSON == "" {
		log.Fatalf("USER_TAGS environment variable is not set")
	}
	err = json.Unmarshal([]byte(userTagsJSON), &userTags)
	if err != nil {
		log.Fatalf("Error parsing USER_TAGS: %v", err)
	}
}

func getDutyPerson(service *calendar.Service) (string, error) {
	t := time.Now().Format(time.RFC3339)
	events, err := service.Events.List(calendarID).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(1).OrderBy("startTime").Do()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve next event: %v", err)
	}

	if len(events.Items) == 0 {
		return "", nil
	}

	return events.Items[0].Summary, nil
}

func sendTelegramMessage(bot *tb.Bot, message string) {
	_, err := bot.Send(&tb.Chat{ID: chatID}, message)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	} else {
		log.Println("Message sent successfully")
	}
}

func main() {
	log.Println("Starting application...")

	log.Printf("TELEGRAM_TOKEN: %s", telegramToken)
	log.Printf("CHAT_ID: %d", chatID)
	log.Printf("CALENDAR_ID: %s", calendarID)
	log.Printf("CRON_EXPRESSION: %s", cronExpression)
	log.Printf("USER_TAGS: %s", userTags)

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Use service account credentials
	config, err := google.JWTConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := config.Client(ctx)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	pref := tb.Settings{
		Token: telegramToken,
	}

	bot, err := tb.NewBot(pref)
	if err != nil {
		log.Fatalf("Unable to create bot: %v", err)
	}

	log.Println("Bot created successfully")

	// Create a time.Location object for the Moscow timezone
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalf("Unable to load location: %v", err)
	}

	// Use local time (Moscow time) for the scheduler
	s := gocron.NewScheduler(location)
	_, err = s.Cron(cronExpression).Do(func() {
		log.Println("Scheduled task started...")
		dutyPerson, err := getDutyPerson(srv)
		if err != nil {
			log.Printf("Error getting duty person: %v", err)
			return
		}

		if dutyPerson != "" {
			tag, exists := userTags[dutyPerson]
			if exists {
				message := fmt.Sprintf("Сегодня дежурит %s %s", dutyPerson, tag)
				log.Println("Preparing to send message...")
				sendTelegramMessage(bot, message)
			} else {
				log.Printf("Duty person %s not found in userTags", dutyPerson)
			}
		} else {
			log.Println("No duty person found")
		}
		log.Println("Scheduled task completed.")
	})
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	s.StartAsync()

	// Keep the program running
	select {}
}
