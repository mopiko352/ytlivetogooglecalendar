package util

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/youtube/v3"
)

// token 取得
func GetConfig(ctx context.Context, secret_path string) (oauth2.TokenSource, error) {

	b, err := GetSecret(ctx, secret_path)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTAccessTokenSourceWithScope(b, youtube.YoutubeScope, calendar.CalendarScope, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}
	return config, nil
}

func GetSecret(ctx context.Context, name string) ([]byte, error) {

	// Create the client.
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %v", err)
	}

	return result.Payload.GetData(), nil
}
