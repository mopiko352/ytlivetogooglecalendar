package app

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"google.golang.org/api/calendar/v3"
)

func TestInsertOrPatchCalenderEvent(t *testing.T) {
	ctx := context.Background()
	secret := os.Getenv("SECRET_PATH")
	if secret == "" {
		t.Fatalf("Set SECRET_PATH environment variables")
	}
	service, err := NewClient(ctx, secret)
	if err != nil {
		t.Fatalf("an error occured %s", err)
	}
	calendarId := os.Getenv("CALENDAR_ID")
	if calendarId == "" {
		t.Fatalf("Set CALENDAR_ID environment variables")
	}
	start := time.Now()
	end := start.Add(1 * time.Hour)
	/*
		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
		b := make([]rune, 10)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		id := string(b)*/
	//New event
	mockevent := &calendar.Event{
		Id:          "abcdef",
		Start:       createEventDateTime(start),
		End:         createEventDateTime(end),
		Summary:     "this is test schedule",
		Description: "TestSchedule",
	}
	log.Printf("%+v", mockevent)
	err = service.InsertOrPatchCalenderEvent(calendarId, mockevent)
	if err != nil {
		t.Fatalf("InsertOrPatchCalenderEvent failed. %s", err)
	}

	//Modify event
	before := &calendar.Event{
		Id:          "test123",
		Start:       createEventDateTime(start),
		End:         createEventDateTime(end),
		Description: "TestSchedule (before modify)",
		Summary:     "this is test schedule (before modify)",
	}
	err = service.InsertOrPatchCalenderEvent(calendarId, before)
	if err != nil {
		t.Fatalf("InsertOrPatchCalenderEvent failed. %s", err)
	}

	end = end.Add(1 * time.Hour)
	after := &calendar.Event{
		Id:          "test123",
		Start:       createEventDateTime(start),
		End:         createEventDateTime(end),
		Description: "TestSchedule (after modify)",
		Summary:     "this is test schedule (after modify)",
	}

	err = service.InsertOrPatchCalenderEvent(calendarId, after)
	if err != nil {
		t.Fatalf("InsertOrPatchCalenderEvent (modify) failed. %s", err)
	}
}
