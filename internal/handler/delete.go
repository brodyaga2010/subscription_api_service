package handler

import (
	"eff_mobile/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Delete godoc
// @Summary      Delete a subscription
// @Description  Delete a subscription by its ID
// @Tags         Subscriptions
// @Param        id  path  int  true  "Subscription ID"
// @Success      204 "No Content"
// @Failure      400 {object} model.Response "Invalid ID format"
// @Failure      500 {object} model.Response "Could not delete subscription"
// @Router       /subscriptions/{id} [delete]
func (api *SubscriptionApi) Delete(c echo.Context) error {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Id is incorrect"})
	}

	if err := api.srvc.DeleteSubscription(c.Request().Context(), idInt); err != nil {
		if err == model.ErrSubscriptionNotFound {
			return echo.NewHTTPError(http.StatusNotFound, model.Response{Message: model.ErrSubscriptionNotFound.Error()})
		}
		api.log.Errorf("Failed delete subscription: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.Response{Message: "Could not delete subscription"})
	}

	return c.NoContent(http.StatusNoContent)
}
