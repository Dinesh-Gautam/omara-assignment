# Deploying the Backend to Google Cloud Run

This guide provides the steps to deploy the backend application to Google Cloud Run.

## Prerequisites

1.  **Google Cloud SDK:** Make sure you have the `gcloud` CLI installed and authenticated.
2.  **Project ID:** Set your Google Cloud Project ID as an environment variable to simplify the commands:
    ```bash
    export PROJECT_ID=$(gcloud config get-value project)
    ```
3.  **Region:** Choose a region for your services (e.g., `us-central1`).
    ```bash
    export REGION=us-central1
    ```

## 1. Enable Required APIs

Enable the necessary Google Cloud APIs for your project.

```bash
gcloud services enable \
  run.googleapis.com \
  sqladmin.googleapis.com \
  cloudbuild.googleapis.com \
  containerregistry.googleapis.com \
  iam.googleapis.com \
  storage.googleapis.com
```

## 2. Create a Cloud SQL for PostgreSQL Instance

```bash
export SQL_INSTANCE_NAME=backend-postgres-instance
export SQL_PASSWORD=$(openssl rand -base64 16) # Generate a strong password

gcloud sql instances create $SQL_INSTANCE_NAME \
  --database-version=POSTGRES_15 \
  --region=$REGION \
  --root-password=$SQL_PASSWORD \
  --project=$PROJECT_ID

# Enable the pgvector extension
gcloud sql instances patch $SQL_INSTANCE_NAME \
    --database-flags=cloudsql.extensions.vector=on

echo "Your Cloud SQL password is: $SQL_PASSWORD"
echo "Please save it securely."
```

## 3. Create a Database and User

```bash
export DB_NAME=strategic_insights
export DB_USER=user

gcloud sql databases create $DB_NAME --instance=$SQL_INSTANCE_NAME --project=$PROJECT_ID
gcloud sql users create $DB_USER --instance=$SQL_INSTANCE_NAME --password=$SQL_PASSWORD --project=$PROJECT_ID
```

## 4. Create Google Cloud Storage Bucket

Create the GCS bucket that the application will use to store uploaded files.

```bash
export GCS_BUCKET_NAME=omara-assignment-bucket
gcloud storage buckets create gs://$GCS_BUCKET_NAME --project=$PROJECT_ID --location=$REGION
```

## 5. Create a Service Account

Create a dedicated service account for your Cloud Run service.

```bash
export SERVICE_ACCOUNT_NAME=cloud-run-backend-sa

gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME \
  --display-name="Cloud Run Backend Service Account" \
  --project=$PROJECT_ID
```

## 6. Grant IAM Roles to the Service Account

Grant the necessary permissions to your new service account.

```bash
# Role for Cloud SQL Client
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

# Role for Firebase Authentication (to get user details)
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/firebaseauth.admin"

# Role for Google Cloud Storage (to read/write to the bucket)
gcloud storage buckets add-iam-policy-binding gs://$GCS_BUCKET_NAME \
  --member="serviceAccount:$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.objectAdmin"
```

## 7. Build the Container Image

Use the optimized Cloud Build configuration to build your Docker image.

```bash
gcloud builds submit --config cloudbuild.yaml .
```

## 8. Deploy to Cloud Run

Deploy your application, connecting it to Cloud SQL and setting the environment variables.

**Important:** Replace `YOUR_GEMINI_API_KEY`, `YOUR_FRONTEND_URL`, and the `SQL_PASSWORD` with your actual values.

```bash
export SQL_CONNECTION_NAME=$(gcloud sql instances describe $SQL_INSTANCE_NAME --format='value(connectionName)')

gcloud run deploy backend \
  --image="gcr.io/$PROJECT_ID/backend:latest" \
  --platform=managed \
  --region=$REGION \
  --service-account="$SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com" \
  --add-cloudsql-instances=$SQL_CONNECTION_NAME \
  --allow-unauthenticated \
  --set-env-vars="POSTGRES_HOST=/cloudsql/$SQL_CONNECTION_NAME" \
  --set-env-vars="POSTGRES_PORT=5432" \
  --set-env-vars="POSTGRES_USER=$DB_USER" \
  --set-env-vars="POSTGRES_PASSWORD=$SQL_PASSWORD" \
  --set-env-vars="POSTGRES_DB=$DB_NAME" \
  --set-env-vars="GCS_BUCKET_NAME=$GCS_BUCKET_NAME" \
  --set-env-vars="GEMINI_API_KEY=YOUR_GEMINI_API_KEY" \
  --set-env-vars="FRONTEND_URL=YOUR_FRONTEND_URL"
```

After the deployment is complete, you will get a URL for your service. Your backend is now deployed!
