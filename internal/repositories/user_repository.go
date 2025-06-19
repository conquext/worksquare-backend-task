package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"housing-api/internal/models"
	"housing-api/internal/utils"
)

// UserRepository handles user data operations
type UserRepository struct {
	users    []models.User
	filePath string
}

// NewUserRepository creates a new user repository
func NewUserRepository() (*UserRepository, error) {
	filePath := utils.GetDataFilePath("users.json")
	repo := &UserRepository{
		filePath: filePath,
		users:    []models.User{},
	}

	// Load existing users if file exists, or create demo user
	if err := repo.loadUsers(); err != nil {
		// If file doesn't exist, create it with demo user
		if err := repo.createDemoUser(); err != nil {
			return nil, fmt.Errorf("failed to create demo user: %w", err)
		}
	}

	return repo, nil
}

// loadUsers loads users from JSON file
func (r *UserRepository) loadUsers() error {
	// Check if file exists
	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		return fmt.Errorf("users file does not exist")
	}

	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read users file: %w", err)
	}

	if len(file) == 0 {
		r.users = []models.User{}
		return nil
	}

	if err := json.Unmarshal(file, &r.users); err != nil {
		return fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return nil
}

// saveUsers saves users to JSON file
func (r *UserRepository) saveUsers() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	data, err := json.MarshalIndent(r.users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	return nil
}

// createDemoUser creates a demo user for testing
func (r *UserRepository) createDemoUser() error {
	// Get demo credentials from environment or use defaults
	demoEmail := os.Getenv("DEMO_USER_EMAIL")
	if demoEmail == "" {
		demoEmail = "demo@worksquare.com"
	}

	demoPassword := os.Getenv("DEMO_USER_PASSWORD")
	if demoPassword == "" {
		demoPassword = "demo123456"
	}

	hashedPassword, err := utils.HashPassword(demoPassword)
	if err != nil {
		return fmt.Errorf("failed to hash demo password: %w", err)
	}

	demoUser := models.User{
		ID:        1,
		Email:     demoEmail,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	r.users = []models.User{demoUser}
	return r.saveUsers()
}

// GetAll returns all users (excluding passwords)
func (r *UserRepository) GetAll() ([]models.User, error) {
	// Return copy without passwords
	var users []models.User
	for _, user := range r.users {
		userCopy := user
		userCopy.Password = "" // Remove password from response
		users = append(users, userCopy)
	}
	return users, nil
}

// GetByID returns a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", id)
}

// GetByEmail returns a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user with email %s not found", email)
}

// Create creates a new user
func (r *UserRepository) Create(user models.User) (*models.User, error) {
	// Check if user with email already exists
	if _, err := r.GetByEmail(user.Email); err == nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Generate new ID
	user.ID = r.getNextID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Add user to slice
	r.users = append(r.users, user)

	// Save to file
	if err := r.saveUsers(); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(id int, updates models.User) (*models.User, error) {
	for i, user := range r.users {
		if user.ID == id {
			// Preserve certain fields
			updates.ID = user.ID
			updates.CreatedAt = user.CreatedAt
			updates.UpdatedAt = time.Now()

			// If password is empty, keep the old password
			if updates.Password == "" {
				updates.Password = user.Password
			}

			// Update the user
			r.users[i] = updates

			// Save to file
			if err := r.saveUsers(); err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}

			return &r.users[i], nil
		}
	}
	return nil, fmt.Errorf("user with ID %d not found", id)
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id int) error {
	for i, user := range r.users {
		if user.ID == id {
			// Remove user from slice
			r.users = append(r.users[:i], r.users[i+1:]...)

			// Save to file
			if err := r.saveUsers(); err != nil {
				return fmt.Errorf("failed to delete user: %w", err)
			}

			return nil
		}
	}
	return fmt.Errorf("user with ID %d not found", id)
}

// EmailExists checks if an email already exists
func (r *UserRepository) EmailExists(email string) bool {
	_, err := r.GetByEmail(email)
	return err == nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id int, newPassword string) error {
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	_, err = r.Update(id, *user)
	return err
}

// GetUserCount returns the total number of users
func (r *UserRepository) GetUserCount() int {
	return len(r.users)
}

// GetRecentUsers returns recently registered users
func (r *UserRepository) GetRecentUsers(limit int) ([]models.User, error) {
	// Sort users by creation date (most recent first)
	sortedUsers := make([]models.User, len(r.users))
	copy(sortedUsers, r.users)

	sort.Slice(sortedUsers, func(i, j int) bool {
		return sortedUsers[i].CreatedAt.After(sortedUsers[j].CreatedAt)
	})

	// Limit results
	if limit > len(sortedUsers) {
		limit = len(sortedUsers)
	}

	result := sortedUsers[:limit]

	// Remove passwords from response
	for i := range result {
		result[i].Password = ""
	}

	return result, nil
}

// SearchUsers searches users by email (partial match)
func (r *UserRepository) SearchUsers(query string) ([]models.User, error) {
	if query == "" {
		return r.GetAll()
	}

	var results []models.User
	query = strings.ToLower(query)

	for _, user := range r.users {
		if strings.Contains(strings.ToLower(user.Email), query) {
			userCopy := user
			userCopy.Password = "" // Remove password
			results = append(results, userCopy)
		}
	}

	return results, nil
}

// ValidateUserCredentials validates user email and password
func (r *UserRepository) ValidateUserCredentials(email, password string) (*models.User, error) {
	user, err := r.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

// UpdateLastLogin updates the user's last login time (if you want to track this)
func (r *UserRepository) UpdateLastLogin(id int) error {
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}

	user.UpdatedAt = time.Now()
	_, err = r.Update(id, *user)
	return err
}

// GetUsersByDateRange returns users created within a date range
func (r *UserRepository) GetUsersByDateRange(start, end time.Time) ([]models.User, error) {
	var results []models.User

	for _, user := range r.users {
		if user.CreatedAt.After(start) && user.CreatedAt.Before(end) {
			userCopy := user
			userCopy.Password = "" // Remove password
			results = append(results, userCopy)
		}
	}

	return results, nil
}

// Backup creates a backup of the users data
func (r *UserRepository) Backup(backupPath string) error {
	data, err := json.MarshalIndent(r.users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users for backup: %w", err)
	}

	return os.WriteFile(backupPath, data, 0644)
}

// Restore restores users data from a backup
func (r *UserRepository) Restore(backupPath string) error {
	file, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(file, &users); err != nil {
		return fmt.Errorf("failed to unmarshal backup data: %w", err)
	}

	r.users = users
	return r.saveUsers()
}

// getNextID generates the next available user ID
func (r *UserRepository) getNextID() int {
	maxID := 0
	for _, user := range r.users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	return maxID + 1
}

// CleanupOldUsers removes users older than specified duration (for maintenance)
func (r *UserRepository) CleanupOldUsers(maxAge time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-maxAge)
	var keptUsers []models.User
	removedCount := 0

	for _, user := range r.users {
		if user.CreatedAt.After(cutoffTime) {
			keptUsers = append(keptUsers, user)
		} else {
			removedCount++
		}
	}

	r.users = keptUsers
	if err := r.saveUsers(); err != nil {
		return 0, fmt.Errorf("failed to save after cleanup: %w", err)
	}

	return removedCount, nil
}