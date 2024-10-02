package endpoints

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"github.com/knackwurstking/picow-led-server/pkg/clients"
)

// TODO: Add paths: "/events/color", "/events/colors"
func createEventsEndpoints(e *echo.Echo, c *clients.Clients) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	e.GET("/events/devices", func(ctx echo.Context) error {
		conn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		client := c.Add(ClientsEmitTypeDevices, conn)

		defer func() {
			c.Remove(ClientsEmitTypeDevices, conn)
		}()

		exit := client.StartHeartBeat()
		defer func() {
			exit <- nil
		}()

		for {
			select {
			case d := <-client.Chan:
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				if err := conn.WriteJSON(d); err != nil {
					return ctx.String(http.StatusInternalServerError, err.Error())
				}
			case <-ctx.Request().Context().Done():
				return ctx.JSON(http.StatusOK, nil)
			case <-client.Done():
				return ctx.JSON(http.StatusOK, nil)
			}
		}
	})

	e.GET("/events/device", func(ctx echo.Context) error {
		conn, err := upgrader.Upgrade(ctx.Response().Writer, ctx.Request(), nil)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}
		client := c.Add(ClientsEmitTypeDevice, conn)

		defer func() {
			c.Remove(ClientsEmitTypeDevice, conn)
		}()

		exit := client.StartHeartBeat()
		defer func() {
			exit <- nil
		}()

		for {
			select {
			case d := <-client.Chan:
				conn.SetWriteDeadline(time.Now().Add(time.Second))
				if err := conn.WriteJSON(d); err != nil {
					return ctx.String(http.StatusInternalServerError, err.Error())
				}
			case <-ctx.Request().Context().Done():
				return ctx.JSON(http.StatusOK, nil)
			case <-client.Done():
				return ctx.JSON(http.StatusOK, nil)
			}
		}
	})
}