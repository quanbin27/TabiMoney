package handlers

import (
    "net/http"
    "strconv"

    "tabimoney/internal/models"
    "tabimoney/internal/services"

    "github.com/labstack/echo/v4"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler { return &CategoryHandler{} }

// List categories: returns system categories and user's categories
func (h *CategoryHandler) List(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    var categories []models.Category
    if err := services.DB().Where("is_system = ? OR user_id = ?", true, userID).Order("sort_order ASC, name ASC").Find(&categories).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to load categories", Message: err.Error()})
    }

    // Map to response
    responses := make([]models.CategoryResponse, 0, len(categories))
    for i := range categories {
        cModel := &categories[i]
        responses = append(responses, models.CategoryResponse{
            ID: cModel.ID,
            UserID: cModel.UserID,
            Name: cModel.Name,
            NameEn: cModel.NameEn,
            Description: cModel.Description,
            ParentID: cModel.ParentID,
            IsSystem: cModel.IsSystem,
            IsActive: cModel.IsActive,
            SortOrder: cModel.SortOrder,
            CreatedAt: cModel.CreatedAt,
            UpdatedAt: cModel.UpdatedAt,
        })
    }
    return c.JSON(http.StatusOK, responses)
}

// Create category: user-defined only
func (h *CategoryHandler) Create(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    var req struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        NameEn string  `json:"name_en"`
        ParentID    *uint64 `json:"parent_id"`
    }
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
    }
    if req.Name == "" {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Validation failed", Message: "name is required"})
    }
    cat := &models.Category{
        UserID:     &userID,
        Name:       req.Name,
        Description:req.Description,
        NameEn:req.NameEn,
        ParentID:   req.ParentID,
        IsSystem:   false,
        IsActive:   true,
    }
    if err := services.DB().Create(cat).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Create failed", Message: err.Error()})
    }
    resp := models.CategoryResponse{
        ID: cat.ID, UserID: cat.UserID, Name: cat.Name, NameEn: cat.NameEn, Description: cat.Description,
        ParentID: cat.ParentID, IsSystem: cat.IsSystem, IsActive: cat.IsActive,
        SortOrder: cat.SortOrder, CreatedAt: cat.CreatedAt, UpdatedAt: cat.UpdatedAt,
    }
    return c.JSON(http.StatusCreated, resp)
}

// Get category by ID
func (h *CategoryHandler) Get(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID", Message: "id must be uint"})
    }

    var cat models.Category
    if err := services.DB().Where("id = ? AND (is_system = ? OR user_id = ?)", id, true, userID).First(&cat).Error; err != nil {
        return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Not found", Message: "Category not found"})
    }

    resp := models.CategoryResponse{
        ID: cat.ID, UserID: cat.UserID, Name: cat.Name, NameEn: cat.NameEn, Description: cat.Description,
        ParentID: cat.ParentID, IsSystem: cat.IsSystem, IsActive: cat.IsActive,
        SortOrder: cat.SortOrder, CreatedAt: cat.CreatedAt, UpdatedAt: cat.UpdatedAt,
    }
    return c.JSON(http.StatusOK, resp)
}

// Update category
func (h *CategoryHandler) Update(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID", Message: "id must be uint"})
    }

    var req struct {
        Name        string  `json:"name"`
        Description string  `json:"description"`
        NameEn string  `json:"name_en"`
        ParentID    *uint64 `json:"parent_id"`
    }
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request", Message: err.Error()})
    }
    if req.Name == "" {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Validation failed", Message: "name is required"})
    }

    var cat models.Category
    if err := services.DB().Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).First(&cat).Error; err != nil {
        return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Not found", Message: "Category not found or cannot be updated"})
    }

    cat.Name = req.Name
    cat.Description = req.Description
    cat.ParentID = req.ParentID
    cat.NameEn = req.NameEn

    if err := services.DB().Save(&cat).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Update failed", Message: err.Error()})
    }

    resp := models.CategoryResponse{
        ID: cat.ID, UserID: cat.UserID, Name: cat.Name, NameEn: cat.NameEn, Description: cat.Description,
        ParentID: cat.ParentID, IsSystem: cat.IsSystem, IsActive: cat.IsActive,
        SortOrder: cat.SortOrder, CreatedAt: cat.CreatedAt, UpdatedAt: cat.UpdatedAt,
    }
    return c.JSON(http.StatusOK, resp)
}

// Delete category
func (h *CategoryHandler) Delete(c echo.Context) error {
    userID := c.Get("user_id").(uint64)
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil {
        return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID", Message: "id must be uint"})
    }

    var cat models.Category
    if err := services.DB().Where("id = ? AND user_id = ? AND is_system = ?", id, userID, false).First(&cat).Error; err != nil {
        return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Not found", Message: "Category not found or cannot be deleted"})
    }

    if err := services.DB().Delete(&cat).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Delete failed", Message: err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}


