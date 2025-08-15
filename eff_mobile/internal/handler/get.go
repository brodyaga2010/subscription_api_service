package handler

import (
	"eff_mobile/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Get godoc
// @Summary      Get a subscription by ID
// @Description  Get details of a specific subscription by its ID
// @Tags         Subscriptions
// @Produce      json
// @Param        id  path      int  true  "Subscription ID"
// @Success      200 {object}  model.Subscription
// @Failure      400 {object}  model.Response "Invalid ID format"
// @Failure      404 {object}  model.Response "Subscription not found"
// @Router       /subscriptions/{id} [get]
func (api *SubscriptionApi) Get(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Id is empty"})
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Id is incorrect"})
	}

	subscription, err := api.srvc.GetSubscription(c.Request().Context(), idInt)
	if err != nil {
		if err == model.ErrSubscriptionNotFound {
			return echo.NewHTTPError(http.StatusNotFound, model.Response{Message: model.ErrSubscriptionNotFound.Error()})
		}
		api.log.Errorf("Failed get subscription: %v", err)

		return echo.NewHTTPError(http.StatusInternalServerError, model.Response{Message: "Could not get subscription"})
	}

	return c.JSON(http.StatusOK, subscription)
}
