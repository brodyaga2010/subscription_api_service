package handler

import (
	"eff_mobile/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// List godoc
// @Summary      List all subscriptions
// @Description  Get a list of all subscriptions
// @Tags         Subscriptions
// @Produce      json
// @Success      200 {array}   model.Subscription
// @Failure      500 {object}  model.Response "Could not list subscriptions"
// @Router       /subscriptions [get]
func (api *SubscriptionApi) List(c echo.Context) error {
	subscriptions, err := api.srvc.ListSubscriptions(c.Request().Context())
	if err != nil {
		if err == model.ErrSubscriptionNotFound {
			return echo.NewHTTPError(http.StatusNotFound, model.Response{Message: model.ErrSubscriptionNotFound.Error()})
		}
		api.log.Errorf("Failed get list subscriptions: %v", err)

		return c.JSON(http.StatusInternalServerError, model.Response{Message: "Could not list subscriptions"})
	}

	return c.JSON(http.StatusOK, subscriptions)
}
