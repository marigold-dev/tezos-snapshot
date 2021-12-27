package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marigold-dev/tezos-snapshot/pkg/util"
	"github.com/patrickmn/go-cache"
)

func main() {
	goCache := cache.New(12*time.Hour, 12*time.Hour)
	bucketName := os.Getenv("BUCKET_NAME")
	e := echo.New()

	e.GET("/items", func(c echo.Context) error {
		items := getSnapshotItemsCached(c.Request().Context(), goCache, bucketName)

		data, err := json.Marshal(items)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, data)
	})

	e.GET("/mainnet", func(c echo.Context) error {
		publicURL, err := getNewestSnapshotsPublicURL(c.Request().Context(), goCache, bucketName, util.MAINNET)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, publicURL)
	})

	e.GET("/testnet", func(c echo.Context) error {
		publicURL, err := getNewestSnapshotsPublicURL(c.Request().Context(), goCache, bucketName, util.TESTNET)
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, publicURL)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
