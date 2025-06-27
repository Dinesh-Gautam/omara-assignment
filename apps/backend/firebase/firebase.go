package firebase

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var AuthClient *auth.Client

func Initialize() {
	ctx := context.Background()
	var app *firebase.App
	var err error

	// When running on Google Cloud (e.g., Cloud Run), GOOGLE_APPLICATION_CREDENTIALS should not be set.
	// The client library will automatically use the service account associated with the resource (ADC).
	// For local development, we use the service account key file.
	credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credsFile != "" {
		opt := option.WithCredentialsFile(credsFile)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// If GOOGLE_APPLICATION_CREDENTIALS is not set, initialize without any options.
		// The SDK will automatically use Application Default Credentials.
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	AuthClient, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
}
