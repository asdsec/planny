package api

import (
	"com.github/asdsec/planny/internal/model"
	"com.github/asdsec/planny/internal/security"
	db "com.github/asdsec/planny/internal/store"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func (serv *Server) createPlan(ctx echo.Context) error {
	var req createPlanRequest
	if err := ctx.Bind(&req); err != nil {
		return serv.err(ctx, http.StatusBadRequest, "cannot bind request body")
	}

	payload := ctx.Get(authorizationPayloadKey).(*security.TokenPayload)
	overlap, err := serv.checkPlanDateOverlap(ctx, payload.UserID, req.StartDate, req.EndDate)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plans not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot check plan date overlap")
	}
	if overlap {
		return serv.err(ctx, http.StatusBadRequest, "plan date overlap")
	}

	arg := db.CreatePlanArg{
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      db.Status(req.Status),
		UserID:      payload.UserID,
	}
	plan, err := serv.store.CreatePlan(ctx.Request().Context(), arg)
	if err != nil {
		return serv.err(ctx, http.StatusInternalServerError, "cannot create plan")
	}

	return ctx.JSON(http.StatusCreated, planResponse(&plan))
}

type (
	createPlanRequest struct {
		Title       string       `json:"title"  validate:"required"`
		Description string       `json:"description" validate:"required"`
		StartDate   time.Time    `json:"start_date" validate:"required"`
		EndDate     time.Time    `json:"end_date" validate:"required"`
		Status      model.Status `json:"status" validate:"required"`
	}

	createPlanResponse struct {
		ID          uint         `json:"id"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		StartDate   time.Time    `json:"start_date"`
		EndDate     time.Time    `json:"end_date"`
		Status      model.Status `json:"status"`
		UserID      uint         `json:"user_id"`
		CreatedAt   time.Time    `json:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at"`
	}
)

func planResponse(plan *model.Plan) *createPlanResponse {
	return &createPlanResponse{
		ID:          plan.ID,
		Title:       plan.Title,
		Description: plan.Description,
		StartDate:   plan.StartDate,
		EndDate:     plan.EndDate,
		Status:      plan.Status,
		UserID:      plan.UserID,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}
}

func (serv *Server) retrievePlans(ctx echo.Context) error {
	payload := ctx.Get(authorizationPayloadKey).(*security.TokenPayload)
	plans, err := serv.store.ListPlansByUserID(ctx.Request().Context(), payload.UserID)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plans not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot retrieve plans")
	}

	return ctx.JSON(http.StatusOK, newRetrievePlansResponse(&plans))
}

type (
	retrievePlanModel createPlanResponse

	retrievePlansResponse struct {
		Plans []retrievePlanModel `json:"plans"`
	}
)

func newRetrievePlansResponse(plans *[]model.Plan) *retrievePlansResponse {
	var res retrievePlansResponse
	for _, plan := range *plans {
		res.Plans = append(res.Plans, retrievePlanModel{
			ID:          plan.ID,
			Title:       plan.Title,
			Description: plan.Description,
			StartDate:   plan.StartDate,
			EndDate:     plan.EndDate,
			Status:      plan.Status,
			UserID:      plan.UserID,
			CreatedAt:   plan.CreatedAt,
			UpdatedAt:   plan.UpdatedAt,
		})
	}
	if len(res.Plans) == 0 {
		res.Plans = []retrievePlanModel{}
	}
	return &res
}

func (serv *Server) deletePlan(ctx echo.Context) error {
	var req deletePlanRequest
	if err := ctx.Bind(&req); err != nil {
		return serv.err(ctx, http.StatusBadRequest, "cannot bind request body")
	}

	payload := ctx.Get(authorizationPayloadKey).(*security.TokenPayload)
	plan, err := serv.store.GetPlanByID(ctx.Request().Context(), req.ID)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plan not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot retrieve plan")
	}
	if plan.UserID != payload.UserID {
		return serv.err(ctx, http.StatusForbidden, "plan does not belong to user")
	}

	err = serv.store.DeletePlanByID(ctx.Request().Context(), req.ID)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plan not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot delete plan")
	}

	return ctx.NoContent(http.StatusNoContent)
}

type (
	deletePlanRequest struct {
		ID uint `param:"id" validate:"required"`
	}
)

func (serv *Server) updatePlan(ctx echo.Context) error {
	var req updatePlanRequest
	if err := ctx.Bind(&req); err != nil {
		return serv.err(ctx, http.StatusBadRequest, "cannot bind request body")
	}

	payload := ctx.Get(authorizationPayloadKey).(*security.TokenPayload)
	plan, err := serv.store.GetPlanByID(ctx.Request().Context(), req.ID)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plan not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot retrieve plan")
	}
	if plan.UserID != payload.UserID {
		return serv.err(ctx, http.StatusForbidden, "plan does not belong to user")
	}

	overlap, err := serv.checkPlanDateOverlap(ctx, payload.UserID, req.StartDate, req.EndDate)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plans not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot check plan date overlap")
	}
	if overlap {
		return serv.err(ctx, http.StatusBadRequest, "plan date overlap")
	}

	arg := db.UpdatePlanArg{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      db.Status(req.Status),
	}
	plan, err = serv.store.UpdatePlanByID(ctx.Request().Context(), arg)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return serv.err(ctx, http.StatusNotFound, "plan not found")
		}
		return serv.err(ctx, http.StatusInternalServerError, "cannot update plan")
	}

	return ctx.JSON(http.StatusOK, planResponse(&plan))
}

type (
	updatePlanRequest struct {
		ID          uint         `param:"id" validate:"required"`
		Title       string       `json:"title"`
		Description string       `json:"description"`
		StartDate   time.Time    `json:"start_date"`
		EndDate     time.Time    `json:"end_date"`
		Status      model.Status `json:"status"`
	}
)

func (serv *Server) checkPlanDateOverlap(ctx echo.Context, userID uint, startDate, endDate time.Time) (bool, error) {
	plans, err := serv.store.ListPlansByUserID(ctx.Request().Context(), userID)
	if err != nil {
		if err.Error() == db.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	for _, plan := range plans {
		if plan.StartDate.Before(endDate) && plan.EndDate.After(startDate) {
			return true, nil
		}
	}

	return false, nil
}
