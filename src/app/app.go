package app

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"

	"github.com/mopiko352/ytlivetogooglecalendar/src/util"

	"github.com/pkg/errors"
	"golang.org/x/xerrors"

	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type client struct {
	context                context.Context
	ChannelsService        *youtube.ChannelsService
	VideosService          *youtube.VideosService
	ChannelSectionsService *youtube.ChannelSectionsService
	SearchService          *youtube.SearchService
	CalenderService        *calendar.Service
}

func NewClient(ctx context.Context, secret_path string) (*client, error) {
	config, err := util.GetConfig(ctx, secret_path)
	if err != nil {
		return nil, err
	}
	ytservice, err := youtube.NewService(ctx, option.WithTokenSource(config))
	if err != nil {
		return nil, errors.Errorf("Unable to create youtube Client %v", err)
	}
	calservice, err := calendar.NewService(ctx, option.WithTokenSource(config))
	if err != nil {
		return nil, errors.Errorf("Unable to create calender Client %v", err)
	}
	channelsService := youtube.NewChannelsService(ytservice)
	videosService := youtube.NewVideosService(ytservice)
	channelsectionsService := youtube.NewChannelSectionsService(ytservice)
	searchService := youtube.NewSearchService(ytservice)

	return &client{
		context:                ctx,
		ChannelsService:        channelsService,
		VideosService:          videosService,
		ChannelSectionsService: channelsectionsService,
		SearchService:          searchService,
		CalenderService:        calservice,
	}, nil
}

func (c *client) searchVideosRequest(id string) *youtube.VideosListCall {
	part := []string{"snippet"}
	return c.VideosService.List(part)
}

func (c *client) SearchVideos(id string) (*youtube.VideoListResponse, error) {
	req := c.searchVideosRequest(id)
	resp, err := req.Do()
	if err != nil {
		return nil, xerrors.Errorf("error when call searchVideos")
	}
	return resp, nil
}

func (c *client) searchChannelRequest(id string) *youtube.ChannelsListCall {
	part := []string{"contentDetails"}
	return c.ChannelsService.List(part).ForUsername(id)
}

func (c *client) SearchChannel(id string) (*youtube.ChannelListResponse, error) {
	req := c.searchChannelRequest(id)
	resp, err := req.Do()
	if err != nil {
		return nil, xerrors.Errorf("error when call searchChannel")
	}
	return resp, nil
}

func (c *client) searchChannelSectionsRequest(id string) *youtube.ChannelSectionsListCall {
	part := []string{"snippet"}
	return c.ChannelSectionsService.List(part).ChannelId(id)
}

func (c *client) SearchChannelSections(id string) (*youtube.ChannelSectionListResponse, error) {
	req := c.searchChannelSectionsRequest(id)
	resp, err := req.Do()
	if err != nil {
		return nil, xerrors.Errorf("error when call searchVideos")
	}
	return resp, nil
}

func (c *client) searchRequest(id string, pageToken string, publishedAfter time.Time) *youtube.SearchListCall {
	part := []string{"snippet"}
	if pageToken != "" {
		return c.SearchService.List(part).ChannelId(id).Type("video").EventType("upcoming").MaxResults(50).Order("date").PageToken(pageToken)
	} else if pageToken != "" && !publishedAfter.IsZero() {
		return c.SearchService.List(part).ChannelId(id).Type("video").EventType("upcoming").MaxResults(50).Order("date").PageToken(pageToken).PublishedAfter(publishedAfter.Format("2006-01-02T15:04:05Z"))
	} else {
		return c.SearchService.List(part).ChannelId(id).Type("video").EventType("upcoming").MaxResults(50).Order("date").PageToken(pageToken)
	}
}

func (c *client) SearchUpcomingLiveAfterDate(id string, date time.Time) ([]*youtube.SearchResult, error) {
	var items []*youtube.SearchResult
	var res func(id string, pageToken string, date time.Time) ([]*youtube.SearchResult, error)
	var end bool
	var pageToken string
	res = func(id string, pageToken string, date time.Time) ([]*youtube.SearchResult, error) {
		req := c.searchRequest(id, pageToken, date)
		resp, err := req.Do()
		if end {
			return items, nil
		}
		if err != nil {
			return nil, xerrors.Errorf("error when call searchVideos:%s", err)
		}
		if resp == nil {
			return nil, xerrors.Errorf("API returns nil")
		}
		for _, item := range resp.Items {
			publishedAt, err := time.Parse("2006-01-02T15:04:05Z", item.Snippet.PublishedAt)
			if err != nil {
				return nil, xerrors.Errorf("youtube data api timeformat something wrong. err:%s format:%s", err, item.Snippet.PublishedAt)
			}
			if publishedAt.After(date) {
				items = append(items, item)
			}
		}
		if pageToken == "" {
			end = true
		}
		pageToken = resp.NextPageToken
		return res(id, pageToken, date)
	}
	return res(id, pageToken, date)
}

func (c *client) ListCalenderAfterDate(calenderId string, date time.Time) ([]*calendar.Event, error) {
	req := c.CalenderService.Events.List(calenderId).MaxResults(50).OrderBy("startTime")
	resp, err := req.Do()
	var res func(calenderId string, date time.Time) ([]*calendar.Event, error)
	var end bool
	var items []*calendar.Event
	res = func(calenderId string, date time.Time) ([]*calendar.Event, error) {
		if end {
			return items, nil
		}
		if err != nil {
			return nil, xerrors.Errorf("error when call ListCalendarEvents:%s", err)
		}
		if resp == nil {
			return nil, xerrors.Errorf("API returns nil")
		}
		for _, item := range resp.Items {
			updatedAt, err := time.Parse("2006-01-02T15:04:05Z", item.Updated)
			if err != nil {
				return nil, xerrors.Errorf("youtube data api timeformat something wrong. err:%s format:%s", err, item.Updated)
			}
			if updatedAt.After(date) {
				items = append(items, item)
			}
		}
		return res(calenderId, date)
	}
	return res(calenderId, date)
}

func createEventDateTime(datetime time.Time) *calendar.EventDateTime {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return &calendar.EventDateTime{
		DateTime: datetime.In(jst).Format(time.RFC3339),
		TimeZone: "Asia/Tokyo",
	}
}

func (c *client) newCalenderEvent(searchresult []*youtube.SearchResult) ([]*calendar.Event, error) {
	var videoids []string
	var events []*calendar.Event
	for _, s := range searchresult {
		videoids = append(videoids, s.Id.VideoId)
	}
	part := []string{"liveStreamingDetails", "snippet"}
	livecontentCall := c.VideosService.List(part).Id(videoids...)
	resp, err := livecontentCall.Do()
	if err != nil {
		return nil, err
	}
	for _, res := range resp.Items {
		start, err := time.Parse("2006-01-02T15:04:05Z", res.LiveStreamingDetails.ScheduledStartTime)
		if err != nil {
			return nil, err
		}
		end := start.Add(90 * time.Minute)
		if err != nil {
			return nil, err
		}
		e := &calendar.Event{
			Id:          fmt.Sprintf("%x", md5.Sum([]byte(res.Id))),
			Location:    fmt.Sprintf("https://www.youtube.com/watch?v=%s", res.Id),
			Start:       createEventDateTime(start),
			End:         createEventDateTime(end),
			Summary:     res.Snippet.Title,
			Description: res.Snippet.Description,
			Visibility:  "public",
			Status:      "confirmed",
		}
		events = append(events, e)
	}
	return events, nil
}

func (c *client) InsertOrPatchCalenderEvent(calendarId string, event *calendar.Event) error {
	getReq := c.CalenderService.Events.Get(calendarId, event.Id)
	_, err := getReq.Do()
	if err != nil {
		insReq := c.CalenderService.Events.Insert(calendarId, event)
		resp, err := insReq.Do()
		if err != nil {
			return xerrors.Errorf("Insert Calender event failed. %s", err)
		}
		log.Printf("Insert calendar title:%s Id:%s", resp.Summary, resp.Id)
	} else {
		patchReq := c.CalenderService.Events.Patch(calendarId, event.Id, event)
		resp, err := patchReq.Do()
		if err != nil {
			return xerrors.Errorf("Patch calendar event failed %s", err)
		}
		log.Printf("Patch calendar title:%s Id:%s", resp.Summary, resp.Id)
	}
	return nil
}

func (c *client) ApplyYTtoCalendar(calendarId string, searchresult []*youtube.SearchResult) error {
	schedules, err := c.newCalenderEvent(searchresult)
	if err != nil {
		return err
	}
	for _, s := range schedules {
		err := c.InsertOrPatchCalenderEvent(calendarId, s)
		if err != nil {
			return err
		}
	}
	return nil
}
