package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mopiko352/ytlivetogooglecalendar/src/app"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/xerrors"
)

func main() {
	// for local.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	_ = os.Setenv("env", "local")
	_ = os.Setenv("FUNCTION_TARGET", "LocalExecute")
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

var secretPath string
var calendarId string
var channelId string

func init() {
	functions.HTTP("LocalExecute", localexecute)
}

func initenv() error {
	secretPath = os.Getenv("SA_SECRET_PATH")
	if secretPath == "" {
		return xerrors.New("Set SECRET_PATH environment variables")
	}
	calendarId = os.Getenv("CALENDAR_ID_SECRET_PATH")
	if calendarId == "" {
		return xerrors.New("Set CALENDAR_ID environment variables")
	}
	channelId = os.Getenv("CHANNEL_ID")
	if calendarId == "" {
		return xerrors.New("Set CHANNEL_ID environment variables")
	}
	return nil
}

// gcfはpub/sub経由で起動するけどこれはローカルで試すようなのでHTTP起動するよ
func localexecute(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := initenv()
	if err != nil {
		http.Error(w, "Set environment variables", 500)
	}
	service, err := app.NewClient(ctx, secretPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured %s", err), 500)
	}
	upcomings, err := service.SearchUpcomingLiveAfterDate(channelId, time.Now().AddDate(-1, 0, 0))
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured %s", err), 500)
	}
	err = service.ApplyYTtoCalendar(calendarId, upcomings)
	if err != nil {
		http.Error(w, fmt.Sprintf("an error occured %s", err), 500)
	}
}
