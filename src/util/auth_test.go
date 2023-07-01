package util

import (
	"context"
	"os"
	"testing"
)

func TestGetSecret(t *testing.T) {
	ctx := context.Background()
	secret := os.Getenv("SECRET_PATH_SA")
	if secret == "" {
		t.Fatalf("Set SECRET_PATH_SA environment variables")
	}
	token, err := GetConfig(ctx, secret)
	if err != nil {
		t.Fatalf("an error occuerd %s", err)
	}
	t.Logf("token:%s", token)

	calendar := os.Getenv("CALENDAR_ID_SECRET_PATH")
	if calendar == "" {
		t.Fatalf("Set CALENDAR_ID_SECRET_PATH environment variables")
	}
	id, err := GetSecret(ctx, calendar)
	if err != nil {
		t.Fatalf("an error occuerd %s", err)
	}
	t.Logf("id:%s", id)
}
