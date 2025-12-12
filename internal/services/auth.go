package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"tabimoney/internal/config"
	"tabimoney/internal/database"
	"tabimoney/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	config *config.Config
	db     *gorm.DB
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		config: cfg,
		db:     database.GetDB(),
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *models.UserCreateRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user with this email or username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var user *models.User
	var profile *models.UserProfile

	// Use transaction to ensure atomicity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user = &models.User{
			Email:        req.Email,
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
			FirstName:    req.FirstName,
			LastName:     req.LastName,
			Phone:        req.Phone,
			IsVerified:   false,
		}

		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Create user profile
		profile = &models.UserProfile{
			UserID:               user.ID,
			MonthlyIncome:        0,
			Currency:             "VND",
			Timezone:             "Asia/Ho_Chi_Minh",
			Language:             "vi",
			NotificationSettings: "{}",
			AISettings:           "{}",
		}

		if err := tx.Create(profile).Error; err != nil {
			return fmt.Errorf("failed to create user profile: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	if err := s.createSession(user.ID, accessToken, refreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(user)

	return &models.AuthResponse{
		User:         s.userToResponse(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(req *models.UserLoginRequest) (*models.AuthResponse, error) {
	// Find user
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	if err := s.createSession(user.ID, accessToken, refreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(user)

	return &models.AuthResponse{
		User:         s.userToResponse(&user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	// Parse refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	// Find user
	var user models.User
	if err := s.db.First(&user, uint64(userID)).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresAt, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update session
	if err := s.updateSession(user.ID, accessToken, newRefreshToken, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &models.AuthResponse{
		User:         s.userToResponse(&user),
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout invalidates user session
func (s *AuthService) Logout(userID uint64, tokenHash string) error {
	// Deactivate session
	if err := s.db.Model(&models.UserSession{}).
		Where("user_id = ? AND token_hash = ?", userID, tokenHash).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}

	// Delete from Redis cache
	ctx := context.Background()
	sessionKey := fmt.Sprintf("session:%s", tokenHash)
	if err := database.DeleteCache(ctx, sessionKey); err != nil {
		// Log error but don't fail logout
		fmt.Printf("Warning: failed to delete session from cache: %v\n", err)
	}

	return nil
}

// ValidateToken validates JWT token and returns user ID
func (s *AuthService) ValidateToken(tokenString string) (uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user ID in token")
	}

	// Allow telegram_access tokens without DB session check
	if t, ok := claims["type"].(string); ok && t == "telegram_access" {
		// Respect exp if present
		if expVal, hasExp := claims["exp"]; hasExp {
			switch exp := expVal.(type) {
			case float64:
				if time.Now().Unix() > int64(exp) {
					return 0, errors.New("token expired")
				}
			case int64:
				if time.Now().Unix() > exp {
					return 0, errors.New("token expired")
				}
			}
		}
		return uint64(userID), nil
	}

	// Default: require active session for normal access tokens
	var session models.UserSession
	if err := s.db.Where("user_id = ? AND token_hash = ? AND is_active = ?",
		uint64(userID), tokenString, true).First(&session).Error; err != nil {
		return 0, errors.New("session not found or inactive")
	}

	if time.Now().After(session.ExpiresAt) {
		return 0, errors.New("token expired")
	}

	return uint64(userID), nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint64, req *models.ChangePasswordRequest) error {
	// Find user
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Deactivate all sessions
	if err := s.db.Model(&models.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error; err != nil {
		return fmt.Errorf("failed to deactivate sessions: %w", err)
	}

	return nil
}

// GetMonthlyIncome returns the user's configured monthly income from profile
func (s *AuthService) GetMonthlyIncome(userID uint64) (float64, error) {
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Ensure profile exists
			profile = models.UserProfile{UserID: userID}
			if err := s.db.Create(&profile).Error; err != nil {
				return 0, fmt.Errorf("failed to create profile: %w", err)
			}
			return profile.MonthlyIncome, nil
		}
		return 0, fmt.Errorf("failed to get profile: %w", err)
	}
	return profile.MonthlyIncome, nil
}

// SetMonthlyIncome updates the user's monthly income in profile
func (s *AuthService) SetMonthlyIncome(userID uint64, amount float64) error {
	if amount < 0 {
		return errors.New("monthly income must be non-negative")
	}
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			profile = models.UserProfile{UserID: userID, MonthlyIncome: amount}
			return s.db.Create(&profile).Error
		}
		return fmt.Errorf("failed to get profile: %w", err)
	}
	profile.MonthlyIncome = amount
	return s.db.Save(&profile).Error
}

// SetLargeTransactionThreshold sets the large transaction threshold for a user
func (s *AuthService) SetLargeTransactionThreshold(userID uint64, threshold *float64) error {
	if threshold != nil && *threshold < 0 {
		return errors.New("large transaction threshold must be non-negative")
	}
	var profile models.UserProfile
	if err := s.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			profile = models.UserProfile{UserID: userID, LargeTransactionThreshold: threshold}
			return s.db.Create(&profile).Error
		}
		return fmt.Errorf("failed to get profile: %w", err)
	}
	profile.LargeTransactionThreshold = threshold
	return s.db.Save(&profile).Error
}

// Telegram Integration Methods

// GenerateTelegramLinkCode generates a link code for Telegram integration
func (s *AuthService) GenerateTelegramLinkCode(userID uint64) (string, error) {
	// Generate random 8-character code
	code := generateRandomCode(8)

	// Store in database with expiration
	linkCode := &models.TelegramLinkCode{
		Code:      code,
		WebUserID: userID,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10 minutes expiry
	}

	if err := s.db.Create(linkCode).Error; err != nil {
		return "", fmt.Errorf("failed to create link code: %w", err)
	}

	return code, nil
}

// ValidateTelegramLinkCode validates a link code and returns web user ID
func (s *AuthService) ValidateTelegramLinkCode(code string) (uint64, error) {
	var linkCode models.TelegramLinkCode

	if err := s.db.Where("code = ? AND expires_at > ? AND used_at IS NULL", code, time.Now()).First(&linkCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("invalid or expired link code")
		}
		return 0, fmt.Errorf("failed to validate link code: %w", err)
	}

	// Mark as used
	now := time.Now()
	linkCode.UsedAt = &now
	if err := s.db.Save(&linkCode).Error; err != nil {
		return 0, fmt.Errorf("failed to mark link code as used: %w", err)
	}

	return linkCode.WebUserID, nil
}

// LinkTelegramAccount links a Telegram user ID with web user ID
func (s *AuthService) LinkTelegramAccount(telegramUserID int64, webUserID uint64) error {
	// Check if already linked
	var existing models.TelegramAccount
	if err := s.db.Where("telegram_user_id = ?", telegramUserID).First(&existing).Error; err == nil {
		// Update existing link
		existing.WebUserID = webUserID
		existing.UpdatedAt = time.Now()
		return s.db.Save(&existing).Error
	}

	// Create new link
	account := &models.TelegramAccount{
		TelegramUserID: telegramUserID,
		WebUserID:      webUserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.db.Create(account).Error
}

// IsTelegramLinked checks if a web user has Telegram linked
func (s *AuthService) IsTelegramLinked(webUserID uint64) (bool, error) {
	var count int64
	if err := s.db.Model(&models.TelegramAccount{}).Where("web_user_id = ?", webUserID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check Telegram link: %w", err)
	}
	return count > 0, nil
}

// DisconnectTelegram unlinks Telegram account for a web user
func (s *AuthService) DisconnectTelegram(webUserID uint64) error {
	return s.db.Where("web_user_id = ?", webUserID).Delete(&models.TelegramAccount{}).Error
}

// GetTelegramUserID gets Telegram user ID for a web user
func (s *AuthService) GetTelegramUserID(webUserID uint64) (int64, error) {
	var account models.TelegramAccount
	if err := s.db.Where("web_user_id = ?", webUserID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("Telegram account not linked")
		}
		return 0, fmt.Errorf("failed to get Telegram user ID: %w", err)
	}
	return account.TelegramUserID, nil
}

// GenerateTelegramJWT generates a permanent JWT token for Telegram user
func (s *AuthService) GenerateTelegramJWT(telegramUserID int64, webUserID uint64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(365 * 24 * time.Hour) // 1 year expiration

	claims := jwt.MapClaims{
		"user_id":          webUserID,
		"telegram_user_id": telegramUserID,
		"type":             "telegram_access",
		"exp":              expiresAt.Unix(),
		"iat":              now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return tokenString, nil
}

// Helper function to generate random code
func generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)

	// Use crypto/rand for secure random generation
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to time-based if crypto/rand fails (shouldn't happen)
		for i := range b {
			b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		}
		return string(b)
	}

	// Map random bytes to charset
	for i := range b {
		b[i] = charset[randomBytes[i]%byte(len(charset))]
	}
	return string(b)
}

// Helper methods

func (s *AuthService) generateTokens(userID uint64) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.GetJWTExpiration())
	refreshExpiresAt := now.Add(s.config.GetJWTRefreshExpiration())

	// Access token
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Refresh token
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     refreshExpiresAt.Unix(),
		"iat":     now.Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, expiresAt, nil
}

func (s *AuthService) createSession(userID uint64, accessToken, refreshToken string, expiresAt time.Time) error {
	session := &models.UserSession{
		UserID:           userID,
		TokenHash:        accessToken,
		RefreshTokenHash: refreshToken,
		ExpiresAt:        expiresAt,
		RefreshExpiresAt: time.Now().Add(s.config.GetJWTRefreshExpiration()),
		IsActive:         true,
	}

	return s.db.Create(session).Error
}

func (s *AuthService) updateSession(userID uint64, accessToken, refreshToken string, expiresAt time.Time) error {
	return s.db.Model(&models.UserSession{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Updates(map[string]interface{}{
			"token_hash":         accessToken,
			"refresh_token_hash": refreshToken,
			"expires_at":         expiresAt,
			"refresh_expires_at": time.Now().Add(s.config.GetJWTRefreshExpiration()),
		}).Error
}

func (s *AuthService) userToResponse(user *models.User) models.UserResponse {
	return UserToResponse(user)
}

// DB exposes the shared database handle for handlers that need direct access
func DB() *gorm.DB {
	return database.GetDB()
}

// UserToResponse converts a User model to API response (exported helper)
func UserToResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Phone:       user.Phone,
		AvatarURL:   user.AvatarURL,
		IsVerified:  user.IsVerified,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
