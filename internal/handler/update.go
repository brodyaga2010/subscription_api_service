package handler

import (
	"eff_mobile/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Update godoc
// @Summary      Update a subscription
// @Description  Update an existing subscription by its ID
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      int                        true  "Subscription ID"
// @Param        subscription  body      model.SubscriptionRequest  true  "Subscription Info to update"
// @Success      200           {object}  model.Subscription
// @Failure      400           {object}  model.Response "Invalid input or ID format"
// @Failure      500           {object}  model.Response "Could not update subscription"
// @Router       /subscriptions/{id} [put]
func (api *SubscriptionApi) Update(c echo.Context) error {
	id := c.Param("id")
	var req model.SubscriptionRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Id is incorrect"})
	}

	subscription, err := api.srvc.UpdateSubscription(c.Request().Context(), idInt, &req)
	if err != nil {
		switch err {
		case model.ErrDateFormat:
			return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: model.ErrDateFormat.Error()})
		case model.ErrSubscriptionNotFound:
			return echo.NewHTTPError(http.StatusNotFound, model.Response{Message: model.ErrSubscriptionNotFound.Error()})
		}
		api.log.Errorf("Failed update subscription: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.Response{Message: "Could not update subscription"})
	}

	return c.JSON(http.StatusOK, subscription)
}
