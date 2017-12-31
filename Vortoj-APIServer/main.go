package main

import (
	_ "strconv"

	"./controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return controllers.RootConnect(c)
	})

	e.GET("/api", func(c echo.Context) error {
		return controllers.RootConnect(c)
	})

	//packet
	e.GET("/api/packet/:id", func(c echo.Context) error {
		return controllers.PacketDataSelectID(c)
	})
	e.GET("/api/packet/new", func(c echo.Context) error {
		return controllers.PacketDataSelectNew(c)
	})
	e.GET("/api/packet/macaddress", func(c echo.Context) error {
		return controllers.PacketDataSelectMacAddress(c)
	})

	//distance
	e.GET("/api/distance/:id", func(c echo.Context) error {
		return controllers.DistanceSelectID(c)
	})
	e.GET("/api/distance/new", func(c echo.Context) error {
		return controllers.DistanceSelectNew(c)
	})
	e.GET("/api/distance/macaddress", func(c echo.Context) error {
		return controllers.DistanceSelectMacAddress(c)
	})
	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
