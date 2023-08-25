# tezos-snapshot

Services to create and provide Tezos nodes snapshots

We have two entry points for two services, the Photographer CronJob and the Web Server.
## Photographer Entrypoint

Service to create Tezos nodes snapshots and upload them into Google Cloud Storage.

We from marigold use it as a CronJob in K8S that is triggered every day.

### How to use

Set the following environment variables:

```bash
export GOOGLE_APPLICATION_CREDENTIALS = "/path/to/your/client_secret.json"
export BUCKET_NAME = "mybucket"
export NETWORK = "MAINNET"
export MAX_DAYS = "3" # optional, default is 7
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

## Server Entrypoint

Service to server Tezos nodes snapshots from Google Cloud and expose them.

We from marigold use it as a Web Service.


### How to use

Set the following environment variables:

```bash
export BUCKET_NAME = "mybucket"
export GOOGLE_APPLICATION_CREDENTIALS = "/path/to/your/client_secret.json"
```

Running locally:

```bash
go run ./cmd/server
```

Running with docker:

```bash
docker build -f server.Dockerfile . -t server
docker run server
```

## Endpoints

* **/** to return json content with all snapshots
* **/mainnet/rolling** to return the last mainnet rolling snapshot
* **/mainnet/full** to return the last mainnet full snapshot
* **/ghostnet/rolling** to return the last testnet rolling snapshot
* **/ghostnet/full** to return the last testnet full snapshot
