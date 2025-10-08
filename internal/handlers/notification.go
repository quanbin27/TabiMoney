package handlers

import (
    "net/http"
    "strconv"
    "tabimoney/internal/services"
    "github.com/labstack/echo/v4"
)

type NotificationHandler struct {
    svc *services.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
    return &NotificationHandler{svc: services.NewNotificationService()}
}

func (h *NotificationHandler) List(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    unread := c.QueryParam("unread") == "true"
    items, err := h.svc.List(userID, unread)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to list notifications", Message: err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]interface{}{"data": items})
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID", Message: "id must be a number"})
    }
    if err := h.svc.MarkRead(userID, id); err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to mark read", Message: err.Error()})
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}


