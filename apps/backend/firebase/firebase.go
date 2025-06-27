package firebase

import (
	"context"
	"log"
	"strategic-insight-analyst/backend/config"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func Initialize() {
	var opt option.ClientOption
	// When running on Google Cloud (e.g., Cloud Run), GOOGLE_APPLICATION_CREDENTIALS is not set.
	// The client library will automatically use the service account associated with the resource.
	// For local development, we use the service account key file.
	if config.AppConfig.GoogleApplicationCredentials != "" {
		opt = option.WithCredentialsFile(config.AppConfig.GoogleApplicationCredentials)
	}

	// If opt is nil, the SDK uses Application Default Credentials.
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	AuthClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
}
