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
	ResponseSize float64 `json:"size"`
	Message      string  `json:"message"`
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func echoPost(c echo.Context) error {
	delayResponse(c.QueryParam("delay"))

	m := &Request{}
	if err := c.Bind(m); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, responseData([]byte(m.Message)))
}

func echoGet(c echo.Context) error {
	delayResponse(c.QueryParam("delay"))
	return c.JSON(http.StatusCreated, responseData(createData(c.QueryParam("size"))))
}

func responseData(responseData []byte) Response {
	response := Response{
		ResponseSize: float64(len(responseData) / 1024),
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
		responseData = make([]byte, size*1024)
		for i := 0; i < size*1024; i++ {
			responseData[i] = 'A'
		}
	}
	return responseData
}

func delayResponse(delayParam string) {
	delay, err := strconv.Atoi(delayParam)
	if err != nil {
		delay = 0
	}
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
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
