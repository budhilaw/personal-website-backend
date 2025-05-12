package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/budhilaw/personal-website-backend/internal/logger"
	"github.com/budhilaw/personal-website-backend/internal/model"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// UserRepository defines methods for user repository
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateProfile(ctx context.Context, id string, user *model.ProfileUpdate) error
	UpdateAvatar(ctx context.Context, id string, avatar string) error
	UpdatePassword(ctx context.Context, id string, password string) error
}

// userRepository is the implementation of UserRepository
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, username, password, email, first_name, last_name, avatar, bio, is_admin, created_at, updated_at 
			  FROM users 
			  WHERE id = $1`

	var user model.User
	var lastName, avatar, bio sql.NullString

	err := r.db.QueryRowxContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.FirstName,
		&lastName,
		&avatar,
		&bio,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		logger.ErrorContext(ctx, "Failed to get user by ID", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	// Set the nullable fields
	if lastName.Valid {
		user.LastName = lastName.String
	}
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}

	return &user, nil
}

// GetByUsername gets a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT id, username, password, email, first_name, last_name, avatar, bio, is_admin, created_at, updated_at 
			  FROM users 
			  WHERE username = $1`

	var user model.User
	var lastName, avatar, bio sql.NullString

	err := r.db.QueryRowxContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.FirstName,
		&lastName,
		&avatar,
		&bio,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		logger.ErrorContext(ctx, "Failed to get user by username", zap.Error(err), zap.String("username", username))
		return nil, err
	}

	// Set the nullable fields
	if lastName.Valid {
		user.LastName = lastName.String
	}
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}

	return &user, nil
}

// GetByEmail gets a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, username, password, email, first_name, last_name, avatar, bio, is_admin, created_at, updated_at 
			  FROM users 
			  WHERE email = $1`

	var user model.User
	var lastName, avatar, bio sql.NullString

	err := r.db.QueryRowxContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.FirstName,
		&lastName,
		&avatar,
		&bio,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		logger.ErrorContext(ctx, "Failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, err
	}

	// Set the nullable fields
	if lastName.Valid {
		user.LastName = lastName.String
	}
	if avatar.Valid {
		user.Avatar = avatar.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}

	return &user, nil
}

// UpdateProfile updates user profile
func (r *userRepository) UpdateProfile(ctx context.Context, id string, profile *model.ProfileUpdate) error {
	query := `UPDATE users 
			  SET first_name = $2, last_name = $3, email = $4, bio = $5, updated_at = $6
			  WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, profile.FirstName, profile.LastName, profile.Email, profile.Bio, time.Now())
	if err != nil {
		logger.ErrorContext(ctx, "Failed to update profile", zap.Error(err), zap.String("id", id))
	}
	return err
}

// UpdateAvatar updates user avatar
func (r *userRepository) UpdateAvatar(ctx context.Context, id string, avatar string) error {
	query := `UPDATE users 
			  SET avatar = $2, updated_at = $3
			  WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, avatar, time.Now())
	if err != nil {
		logger.ErrorContext(ctx, "Failed to update avatar", zap.Error(err), zap.String("id", id))
	}
	return err
}

// UpdatePassword updates user password
func (r *userRepository) UpdatePassword(ctx context.Context, id string, password string) error {
	query := `UPDATE users 
			  SET password = $2, updated_at = $3
			  WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, password, time.Now())
	if err != nil {
		logger.ErrorContext(ctx, "Failed to update password", zap.Error(err), zap.String("id", id))
	}
	return err
}
