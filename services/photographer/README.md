# Photographer

Service to create Tezos nodes snapshots and upload them into Google Cloud Storage.

We from marigold use it as a CronJob in K8S that is triggered every day.

## How to use

Set the following environment variables:

```bash
export BUCKET_NAME = "mybucket"
export MAX_DAYS = "3" # optional, default is 7
export GOOGLE_APPLICATION_CREDENTIALS = "/path/to/your/client_secret.json"
```

Running locally:

```bash
go run /services/photographer
```

Running with docker:

```bash
docker build -f photographer.Dockerfile . -t photographer
docker run photographer
```