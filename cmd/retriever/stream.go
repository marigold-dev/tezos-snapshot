package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func streamFile(e echo.Context, client *http.Client, fileName, url string) error {
	//Copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	w := e.Response().Writer

	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(fileName))
	w.Header().Set("Content-Type", e.Response().Header().Get("Content-Type"))
	w.Header().Set("Content-Length", e.Response().Header().Get("Content-Length"))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, resp.Body)
	return err
}
