package user

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/pkg/middleware"
	"brief/pkg/repository/storage/postgres"
	"brief/utility"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Register contains business logic for registering a new user
func Register(user *model.User, isAdmin ...bool) (string, error) {
	// Hash password
	hash, salt, err := utility.HashPassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("could not hash user password, got error: %w", err)
	}

	// Set/Modify some fields from the request struct
	{
		user.ID = uuid.New().String()
		user.Password = hash
		user.Salt = salt
		user.CreatedAt = time.Now()
		if (len(isAdmin) > 0) && isAdmin[0] {
			user.Role = constant.Roles[constant.Admin]
		} else {
			user.Role = constant.Roles[constant.User]
		}
	}

	token, err := middleware.CreateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", fmt.Errorf("could not create token, got error: %w", err)
	}

	db := postgres.GetDB()
	err = db.CreateUser(context.TODO(), user)
	if err != nil {
		return "", fmt.Errorf("could not create user, got error: %w", err)
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return token, nil
}

// Login contains business logic for logging in
func Login(userLogin *model.UserLogin) (*model.LoginResponse, error) {
	db := postgres.GetDB()
	user, err := db.GetUser(context.TODO(), userLogin.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("could not get user, got error %w", err)
		}
		return nil, fmt.Errorf("invalid user")
	}

	if !utility.PasswordIsValid(userLogin.Password, user.Salt, user.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	// Ensure that user is not locked
	if user.IsLocked {
		return nil, fmt.Errorf("cannot login, user is currently locked")
	}

	token, err := middleware.CreateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("could not create token")
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return &model.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Get contains business logic to get a user by id or email
func Get(idOrEmail string) (*model.User, error) {
	db := postgres.GetDB()
	user, err := db.GetUser(context.TODO(), idOrEmail)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("could not get user, got error %w", err)
		}
		return nil, fmt.Errorf("invalid user")
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return user, nil
}

// GetAll contains business logic for fetching all users
func GetAll() ([]model.User, error) {
	db := postgres.GetDB()
	users, err := db.GetAllUsers(context.TODO())
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("could not get user, got error: %w", err)
		}
		return nil, fmt.Errorf("invalid user")
	}

	return users, nil
}

// Update contains business logic to update a user
func Update(id string, user *model.User) error {
	db := postgres.GetDB()

	fUser, err := db.GetUser(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Ensure that user is not locked
	if fUser.IsLocked {
		return fmt.Errorf("cannot update, user is currently locked")
	}

	if err := db.UpdateUser(context.TODO(), id, user); err != nil {
		return fmt.Errorf("could not update user, got error: %w", err)
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return nil
}

// ResetPassword contains business logic to reset a user's password
func ResetPassword(id string, rp *model.ResetPassword) (*model.User, error) {
	db := postgres.GetDB()
	fUser, err := db.GetUser(context.TODO(), id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Ensure that user is not locked
	if fUser.IsLocked {
		return nil, fmt.Errorf("cannot update, user is currently locked")
	}

	hashedPassword, salt, err := utility.HashPassword(rp.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password, got error: %w", err)
	}

	// Set hashed password and salt
	rp.Password = hashedPassword
	rp.Salt = salt

	user, err := db.ResetPassword(context.TODO(), id, rp)
	if err != nil {
		return nil, fmt.Errorf("could not reset password, got error: %w", err)
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return user, nil
}

// ForgotPassword contains business logic to handle a forgot-password request
func ForgotPassword(email *model.ForgotPassword) error {
	db := postgres.GetDB()
	user, err := db.GetUser(context.TODO(), email.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return fmt.Errorf("could not fetch user, got error: %w", err)
		}
		return fmt.Errorf("user does not exist")
	}

	_ = user // Prevents the Go compiler from calling an error

	// ### Uncomment the lines below to retrieve token to be sent in email ###

	//token, err := middleware.CreateToken(user.ID, user.Email, user.Role)
	//if err != nil {
	//	return fmt.Errorf("could not create token, got error: %w", err)
	//}

	// TODO: Send Forgot-Password Email

	return nil
}

// LockUser contains business logic to lock a user's account
func LockUser(idOrEmail string) (*model.User, error) {
	db := postgres.GetDB()

	user, err := db.LockUnlock(context.TODO(), idOrEmail, true)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("could not unlock user, got error: %w", err)
		}
		return nil, fmt.Errorf("user does not exist")
	}

	return user, nil
}

// UnlockUser contains business logic to unlock a user's account
func UnlockUser(idOrEmail string) (*model.User, error) {
	db := postgres.GetDB()

	user, err := db.LockUnlock(context.TODO(), idOrEmail, false)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("could not unlock user, got error: %w", err)
		}
		return nil, fmt.Errorf("user does not exist")
	}

	return user, nil
}
