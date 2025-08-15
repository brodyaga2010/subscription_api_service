package handler

import (
	"eff_mobile/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create a new subscription
// @Description  Adds a new subscription to the database
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      model.SubscriptionRequest  true  "Subscription Info"
// @Success      201           {object}  model.ResponseID            "Returns the ID of the created subscription"
// @Failure      400           {object}  model.Response             "Invalid input or date format"
// @Failure      500           {object}  model.Response             "Could not create subscription"
// @Router       /subscriptions [post]
func (api *SubscriptionApi) Create(c echo.Context) error {
	var req model.SubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Invalid input"})
	}
	subscription, err := api.srvc.CreateSubscription(c.Request().Context(), &req)
	if err != nil {
		switch err {
		case model.ErrDateFormat:
			return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: model.ErrDateFormat.Error()})
		case model.ErrCreateSubscription:
			return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: model.ErrCreateSubscription.Error()})
		}
		api.log.Errorf("Failed create row: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.Response{Message: "Could not create subscription"})
	}

	return c.JSON(http.StatusCreated, model.ResponseID{ID: subscription.ID})
}
