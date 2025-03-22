package logic

import (
	"context"
	"log"
	"time"

	"github.com/overcout/Inferno-AI/internal/store"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// CreateEventGoogle executes creation of a calendar event using the user's token
func CreateEventGoogle(user *store.User, cmd *CreateEventGoogleCommand) error {
	ctx := context.Background()

	service, err := calendar.NewService(ctx, option.WithTokenSource(&store.UserTokenSource{User: user}))
	if err != nil {
		log.Println("[GOOGLE] Failed to init calendar client:", err)
		return err
	}

	startTime, err := time.Parse("2006-01-02T15:04:05", cmd.StartTime)
	if err != nil {
		log.Println("[GOOGLE] Invalid start time format:", err)
		return err
	}

	endTime := startTime.Add(time.Duration(cmd.DurationMinutes) * time.Minute)

	event := &calendar.Event{
		Summary: cmd.Title,
		Start:   &calendar.EventDateTime{DateTime: startTime.Format(time.RFC3339)},
		End:     &calendar.EventDateTime{DateTime: endTime.Format(time.RFC3339)},
	}

	_, err = service.Events.Insert("primary", event).Do()
	if err != nil {
		log.Println("[GOOGLE] Failed to insert event:", err)
		return err
	}

	log.Println("[GOOGLE] Event created successfully")
	return nil
}