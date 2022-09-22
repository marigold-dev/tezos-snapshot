package main

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/patrickmn/go-cache"
)

func main() {
	goCache := cache.New(1*time.Hour, 1*time.Hour)
	bucketName := os.Getenv("BUCKET_NAME")
	timeout := time.Duration(5) * time.Second
	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	downloadableHandlerBuilder := func(network snapshot.NetworkType, protocol snapshot.NetworkProtocolType) func(c echo.Context) error {
		return func(c echo.Context) error {
			snapshotType := snapshot.ROLLING
			if c.Param("type") == "full" {
				snapshotType = snapshot.FULL
			}

			snapshot, err := getNewestSnapshot(c.Request().Context(), goCache, bucketName, network, snapshotType, protocol)
			if err != nil {
				return err
			}

			return streamFile(c, client, snapshot.FileName, snapshot.PublicURL)
		}
	}

	e.GET("/mainnet", downloadableHandlerBuilder(snapshot.MAINNET, snapshot.MAIN))
	e.GET("/mainnet/:type", downloadableHandlerBuilder(snapshot.MAINNET, snapshot.MAIN))
	e.GET("/testnet", downloadableHandlerBuilder(snapshot.TESTNET, snapshot.JAKARTA))
	e.GET("/testnet/:type", downloadableHandlerBuilder(snapshot.TESTNET, snapshot.JAKARTA))
	e.GET("/hangzhounet/:type", downloadableHandlerBuilder(snapshot.TESTNET, snapshot.HANGZHOU))
	e.GET("/ithacanet/:type", downloadableHandlerBuilder(snapshot.TESTNET, snapshot.ITHACA))
	e.GET("/jakartanet/:type", downloadableHandlerBuilder(snapshot.TESTNET, snapshot.JAKARTA))
	e.GET("/kathmandunet/:type", downloadableHandlerBuilder(snapshot.TESTNET, snapshot.KATHMANDUNET))
	e.GET("/", func(c echo.Context) error {
		snapshots := getSnapshotItemsCached(c.Request().Context(), goCache, bucketName)
		responseSnapshot := []snapshot.SnapshotItem{}

		for _, i := range snapshots {
			responseSnapshot = append(responseSnapshot, i)
		}

		return c.JSON(http.StatusOK, responseSnapshot)
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "UP")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
