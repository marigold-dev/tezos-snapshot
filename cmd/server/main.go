package main

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/marigold-dev/tezos-snapshot/pkg/snapshot"
	"github.com/patrickmn/go-cache"
)

func main() {
	// godotenv.Load("../../.env")

	goCache := cache.New(5*time.Minute, 10*time.Minute)
	bucketName := os.Getenv("BUCKET_NAME")
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	downloadableHandlerBuilder := func(chain string) func(c echo.Context) error {
		return func(c echo.Context) error {
			historyMode := snapshot.ROLLING
			if c.Param("type") == "full" {
				historyMode = snapshot.FULL
			}

			snapshot, err := getNewestSnapshot(c.Request().Context(), goCache, bucketName, historyMode, chain)
			if err != nil {
				return err
			}

			return c.Redirect(http.StatusFound, snapshot.URL)
		}
	}
	api := func(c echo.Context) error {
		responseCached := getSnapshotResponseCached(c.Request().Context(), goCache, bucketName)
		return c.JSON(http.StatusOK, &responseCached)
	}

	e.GET("/mainnet", downloadableHandlerBuilder("mainnet"))
	e.GET("/mainnet/:type", downloadableHandlerBuilder("mainnet"))
	e.GET("/hangzhounet/:type", downloadableHandlerBuilder("hangzhounet"))
	e.GET("/ghostnet/:type", downloadableHandlerBuilder("ghostnet"))
	e.GET("/ithacanet/:type", downloadableHandlerBuilder("ithacanet"))
	e.GET("/jakartanet/:type", downloadableHandlerBuilder("jakarta"))
	e.GET("/kathmandunet/:type", downloadableHandlerBuilder("kathmandunet"))
	e.GET("/limanet/:type", downloadableHandlerBuilder("limanet"))
	e.GET("/mumbainet/:type", downloadableHandlerBuilder("mumbainet"))
	e.GET("/nairobinet/:type", downloadableHandlerBuilder("nairobinet"))
	e.GET("/oxfordnet/:type", downloadableHandlerBuilder("oxfordnet"))
	e.GET("/", api)
	e.GET("/tezos-snapshots.json", api)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "UP")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type SnapshotResponse struct {
	DateGenerated string                  `json:"date_generated"`
	Org           string                  `json:"org"`
	Schema        string                  `json:"$schema"`
	Data          []snapshot.SnapshotItem `json:"data"`
}
