package services

import (
	"fmt"
	"time"

	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"gorm.io/gorm"
)

type GoalService struct {
	db *gorm.DB
}

func NewGoalService() *GoalService {
	return &GoalService{
		db: database.GetDB(),
	}
}

// CreateGoal creates a new financial goal
func (s *GoalService) CreateGoal(userID uint64, req *models.FinancialGoalCreateRequest) (*models.FinancialGoal, error) {
	goal := &models.FinancialGoal{
		UserID:       userID,
		Title:        req.Title,
		Description:  req.Description,
		TargetAmount: req.TargetAmount,
		CurrentAmount: 0,
		TargetDate:   req.TargetDate,
		GoalType:     req.GoalType,
		Priority:     req.Priority,
		IsAchieved:   false,
	}

	if err := s.db.Create(goal).Error; err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	return goal, nil
}

// GetGoals retrieves user's financial goals
func (s *GoalService) GetGoals(userID uint64) ([]models.FinancialGoal, error) {
	var goals []models.FinancialGoal
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error; err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}

	// Calculate progress for each goal
	for i := range goals {
		if goals[i].TargetAmount <= 0 {
			goals[i].Progress = 0
		} else {
			goals[i].Progress = (goals[i].CurrentAmount / goals[i].TargetAmount) * 100
		}
	}

	return goals, nil
}

// UpdateGoal updates an existing goal
func (s *GoalService) UpdateGoal(userID uint64, goalID uint64, req *models.FinancialGoalUpdateRequest) (*models.FinancialGoal, error) {
	var goal models.FinancialGoal
	if err := s.db.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		return nil, fmt.Errorf("goal not found: %w", err)
	}

	// Update fields
	goal.Title = req.Title
	goal.Description = req.Description
	goal.TargetAmount = req.TargetAmount
	goal.CurrentAmount = req.CurrentAmount
	goal.TargetDate = req.TargetDate
	goal.GoalType = req.GoalType
	goal.Priority = req.Priority

	// Check if goal is achieved
	if goal.CurrentAmount >= goal.TargetAmount && !goal.IsAchieved {
		goal.IsAchieved = true
		now := time.Now()
		goal.AchievedAt = &now
	}

	if err := s.db.Save(&goal).Error; err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	// Calculate progress
	if goal.TargetAmount <= 0 {
		goal.Progress = 0
	} else {
		goal.Progress = (goal.CurrentAmount / goal.TargetAmount) * 100
	}

	return &goal, nil
}

// DeleteGoal deletes a goal
func (s *GoalService) DeleteGoal(userID uint64, goalID uint64) error {
	result := s.db.Where("id = ? AND user_id = ?", goalID, userID).Delete(&models.FinancialGoal{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete goal: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("goal not found")
	}
	return nil
}

// AddContribution adds money to a goal
func (s *GoalService) AddContribution(userID uint64, goalID uint64, amount float64, note string) (*models.FinancialGoal, error) {
	var goal models.FinancialGoal
	if err := s.db.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		return nil, fmt.Errorf("goal not found: %w", err)
	}

	goal.CurrentAmount += amount

	// Check if goal is achieved
	if goal.CurrentAmount >= goal.TargetAmount && !goal.IsAchieved {
		goal.IsAchieved = true
		now := time.Now()
		goal.AchievedAt = &now
	}

	if err := s.db.Save(&goal).Error; err != nil {
		return nil, fmt.Errorf("failed to add contribution: %w", err)
	}

	// Calculate progress
	if goal.TargetAmount <= 0 {
		goal.Progress = 0
	} else {
		goal.Progress = (goal.CurrentAmount / goal.TargetAmount) * 100
	}

	return &goal, nil
}
