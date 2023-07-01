package function

import (
	"context"
	"os"
	"time"

	"github.com/mopiko352/ytlivetogooglecalendar/src/app"
	"github.com/mopiko352/ytlivetogooglecalendar/src/util"

	"golang.org/x/xerrors"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

type MessagePublishedData struct {
	Message PubSubMessage
}

type PubSubMessage struct {
	Data []byte `json:"data"`
}

var secretPath string
var calendarPath string
var channelId string

func init() {
	functions.CloudEvent("LiveToCalendar", LiveToCalendar)
}

func initenv() error {
	secretPath = os.Getenv("SA_SECRET_PATH")
	if secretPath == "" {
		return xerrors.New("Set SECRET_PATH environment variables")
	}
	calendarPath = os.Getenv("CALENDAR_ID_SECRET_PATH")
	if calendarPath == "" {
		return xerrors.New("Set CALENDAR_ID environment variables")
	}
	channelId = os.Getenv("CHANNEL_ID")
	if channelId == "" {
		return xerrors.New("Set CALENDAR_ID environment variables")
	}
	return nil
}

func LiveToCalendar(ctx context.Context, e event.Event) error {
	//TODO: receive parameter from pubsub
	err := initenv()
	if err != nil {
		return xerrors.Errorf("error when init env: %s", err)
	}
	calendarId, err := util.GetSecret(ctx, calendarPath)
	if err != nil {
		return xerrors.Errorf("error when GetcalendarId: %s", err)
	}
	service, err := app.NewClient(ctx, secretPath)
	if err != nil {
		return xerrors.Errorf("error when create client: %s", err)
	}
	upcomings, err := service.SearchUpcomingLiveAfterDate(channelId, time.Now().AddDate(-1, 0, 0))
	if err != nil {
		return xerrors.Errorf("an error occured %s", err)
	}
	err = service.ApplyYTtoCalendar(string(calendarId), upcomings)
	if err != nil {
		return xerrors.Errorf("error when apply calendar %s", err)
	}

	return nil
}
