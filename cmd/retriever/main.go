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
	goCache := cache.New(12*time.Hour, 12*time.Hour)
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

	downloadableHandlerBuilder := func(network snapshot.NetworkType) func(c echo.Context) error {
		return func(c echo.Context) error {
			snapshotType := snapshot.ROLLING
			if c.Param("type") == "full" {
				snapshotType = snapshot.FULL
			}

			snapshot, err := getNewestSnapshot(c.Request().Context(), goCache, bucketName, network, snapshotType)
			if err != nil {
				return err
			}

			return streamFile(c, client, snapshot.FileName, snapshot.PublicURL)
		}
	}

	e.GET("/mainnet", downloadableHandlerBuilder(snapshot.MAINNET))
	e.GET("/mainnet/:type", downloadableHandlerBuilder(snapshot.MAINNET))
	e.GET("/testnet", downloadableHandlerBuilder(snapshot.TESTNET))
	e.GET("/testnet/:type", downloadableHandlerBuilder(snapshot.TESTNET))
	e.GET("/", func(c echo.Context) error {
		items := getSnapshotItemsCached(c.Request().Context(), goCache, bucketName)
		return c.JSON(http.StatusOK, items)
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "UP")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
