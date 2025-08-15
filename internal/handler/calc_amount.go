package handler

import (
	"eff_mobile/internal/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CalculateAmount godoc
// @Summary      Calculate total subscription cost
// @Description  Calculates the sum of subscription prices for a given period, with optional filters for user_id and service_name.
// @Tags         Subscriptions
// @Produce      json
// @Param        from          query     string  true   "Start period (MM-YYYY)"  example("07-2025")
// @Param        to            query     string  true   "End period (MM-YYYY)"    example("08-2025")
// @Param        user_id       query     string  false  "User ID (UUID)"          example("60601fee-2bf1-4721-ae6f-7636e79a0cba")
// @Param        service_name  query     string  false  "Service Name"            example("Yandex Plus")
// @Success      200           {object}  model.SumResponse
// @Failure      400           {object}  model.Response "Invalid query params or date format"
// @Failure      500           {object}  model.Response "Failed to calculate amount"
// @Router       /subscriptions/sum [get]
func (api *SubscriptionApi) CalculateAmount(c echo.Context) error {
	var req model.SumRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Invalid query params"})
	}

	total, err := api.srvc.CalculateAmount(c.Request().Context(), req)
	if err != nil {
		switch err {
		case model.ErrDateIsNull:
			return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: model.ErrDateIsNull.Error()})
		case model.ErrDateFormat:
			return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: model.ErrDateFormat.Error()})
		case model.ErrDateBefore:
			return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: model.ErrDateBefore.Error()})
		}
		api.log.Error("Failed to calculate amount:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, model.Response{Message: "Failed to calculate amount"})
	}

	return c.JSON(http.StatusOK, total)
}
