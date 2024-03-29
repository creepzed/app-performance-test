package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strconv"
	"time"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	ResponseSize int    `json:"size"`
	Message      string `json:"message"`
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func echoPost(c echo.Context) error {
	m := &Request{}
	if err := c.Bind(m); err != nil {
		return err
	}
	response := responseData([]byte(m.Message))
	time.Sleep(responseDelay(c.QueryParam("delay")))
	return c.JSON(http.StatusCreated, response)
}

func echoGet(c echo.Context) error {
	response := responseData(createData(c.QueryParam("size")))
	time.Sleep(responseDelay(c.QueryParam("delay")))
	return c.JSON(http.StatusOK, response)
}

func responseData(responseData []byte) Response {
	response := Response{
		ResponseSize: len(responseData),
		Message:      string(responseData),
	}
	return response
}

func createData(sizeParam string) []byte {
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		size = 0
	}
	var responseData []byte
	if size > 0 {
		responseData = make([]byte, size)
		for i := 0; i < size; i++ {
			responseData[i] = 'A'
		}
	}
	return responseData
}

func responseDelay(delayParam string) time.Duration {
	delay, err := strconv.Atoi(delayParam)
	if err != nil {
		delay = 0
	}
	return time.Duration(delay) * time.Millisecond
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.HideBanner = true

	e.GET("/ping", ping)

	v1 := e.Group("/v1")
	{
		subGroup := v1.Group("/echo")
		{
			subGroup.POST("", echoPost)
			subGroup.GET("", echoGet)
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
