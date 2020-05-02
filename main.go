package main

import (
	"github.com/ahmetask/gocache/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type CustomCache struct {
	Key   string      `json:"key"`
	Life  int64       `json:"life"`
	Value interface{} `json:"value"`
}

type SaveCachingMiddleWare struct {
	service string
}

type GetCacheMiddleWare struct {
	service string
}

func (m *GetCacheMiddleWare) GetCacheMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != "GET" {
				nextError := next(c)
				return nextError
			}

			key := c.Param("key")

			var r interface{}
			err := gocache.GetCache(m.service, key, &r)
			if err != nil {
				c.Set("cache-error", err)
			} else {
				return c.JSON(http.StatusOK, r)
			}

			nextError := next(c)
			return nextError
		}
	}
}

func (m *SaveCachingMiddleWare) SaveCacheMiddleware() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != "POST" {
				nextError := next(c)
				return nextError
			}

			ch := new(CustomCache)
			if err := c.Bind(ch); err != nil {
				return err
			}

			request := gocache.AddCacheRequest{Key: ch.Key, Value: ch.Value, Life: ch.Life}
			err := gocache.SaveCache(m.service, request)
			if err != nil {
				c.Set("cache-error", err)
			}

			nextError := next(c)
			return nextError
		}
	}
}

func main() {
	e := echo.New()

	//serviceName := "localhost:10001"
	serviceName := "gocache-server:10001"

	saveCacheMiddleWare := SaveCachingMiddleWare{
		service: serviceName,
	}

	getCacheMiddleWare := GetCacheMiddleWare{
		service: serviceName,
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(saveCacheMiddleWare.SaveCacheMiddleware())
	e.Use(getCacheMiddleWare.GetCacheMiddleware())

	e.GET("/cache/:key", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "service-response")
	})

	e.POST("/cache", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("cache-error"))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
