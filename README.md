# tezos-snapshot

Services to create and provide Tezos nodes snapshots

We have two entry points for two services, the Photographer CronJob and the Retriever web server.
## Photographer Entrypoint

Service to create Tezos nodes snapshots and upload them into Google Cloud Storage.

We from marigold use it as a CronJob in K8S that is triggered every day.

### How to use

Set the following environment variables:

```bash
export BUCKET_NAME = "mybucket"
export MAX_DAYS = "3" # optional, default is 7
export GOOGLE_APPLICATION_CREDENTIALS = "/path/to/your/client_secret.json"
```

Running locally:

```bash
go run ./cmd/photographer
```

Running with docker:

```bash
docker build -f photographer.Dockerfile . -t photographer
docker run photographer
```

## Retriever Entrypoint

Service to retriever Tezos nodes snapshots from Google Cloud and expose them.

We from marigold use it as a Web Service.


### How to use

Set the following environment variables:

```bash
export BUCKET_NAME = "mybucket"
export GOOGLE_APPLICATION_CREDENTIALS = "/path/to/your/client_secret.json"
```

Running locally:

```bash
go run ./cmd/retriever
```

Running with docker:

```bash
docker build -f retriever.Dockerfile . -t retriever
docker run retriever
```

## Endpoints

* **/items** to return json content with all snapshots
* **/mainnet** to return mainnet downloadable publicURL
* **/testnet** to return testnet downloadable publicURL