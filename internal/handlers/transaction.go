package handlers

import (
    "net/http"
    "strconv"
    "time"

    "tabimoney/internal/config"
    "tabimoney/internal/models"
    "tabimoney/internal/services"

    "github.com/labstack/echo/v4"
)

type TransactionHandler struct {
    svc *services.TransactionService
}

func NewTransactionHandler(cfg *config.Config) *TransactionHandler {
    return &TransactionHandler{svc: services.NewTransactionService(cfg)}
}

func (h *TransactionHandler) List(c echo.Context) error {
    userID := c.Get("user_id").(uint64)

    page, _ := strconv.Atoi(c.QueryParam("page"))
    if page <= 0 { page = 1 }
    limit, _ := strconv.Atoi(c.QueryParam("limit"))
    if limit <= 0 { limit = 20 }

    var categoryID *uint64
    if v := c.QueryParam("category_id"); v != "" {
        if id64, err := strconv.ParseUint(v, 10, 64); err == nil {
            cid := uint64(id64)
            categoryID = &cid
        }
    }

    var txType *string
    if v := c.QueryParam("transaction_type"); v != "" { txType = &v }

    var startDate *time.Time
    if v := c.QueryParam("start_date"); v != "" {
        if t, err := time.Parse("2006-01-02", v); err == nil { startDate = &t }
    }
    var endDate *time.Time
    if v := c.QueryParam("end_date"); v != "" {
        if t, err := time.Parse("2006-01-02", v); err == nil { endDate = &t }
    }

    var minAmount *float64
    if v := c.QueryParam("min_amount"); v != "" {
        if f, err := strconv.ParseFloat(v, 64); err == nil { minAmount = &f }
    }
    var maxAmount *float64
    if v := c.QueryParam("max_amount"); v != "" {
        if f, err := strconv.ParseFloat(v, 64); err == nil { maxAmount = &f }
    }

    req := &models.TransactionQueryRequest{
        Page: page,
        Limit: limit,
        CategoryID: categoryID,
        TransactionType: txType,
        StartDate: startDate,
        EndDate: endDate,
        MinAmount: minAmount,
        MaxAmount: maxAmount,
        Search: c.QueryParam("search"),
        SortBy: c.QueryParam("sort_by"),
        SortOrder: c.QueryParam("sort_order"),
    }

    items, total, err := h.svc.GetTransactions(userID, req)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to list transactions", Message: err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "data": items,
        "total": total,
        "page": page,
        "limit": limit,
    })
}

func (h *TransactionHandler) Create(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    var req models.TransactionCreateRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
    }
    tx, err := h.svc.CreateTransaction(userID, &req)
    if err != nil { return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Create failed", Message: err.Error()}) }
    return c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) Update(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil { return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID", Message: "id must be uint"}) }
    var req models.TransactionUpdateRequest
    if err := c.Bind(&req); err != nil { return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()}) }
    tx, err := h.svc.UpdateTransaction(userID, uint64(id), &req)
    if err != nil { return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Update failed", Message: err.Error()}) }
    return c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) Delete(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil { return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID", Message: "id must be uint"}) }
    if err := h.svc.DeleteTransaction(userID, uint64(id)); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Delete failed", Message: err.Error()})
    }
    return c.JSON(http.StatusOK, SuccessResponse{Message: "Deleted"})
}


