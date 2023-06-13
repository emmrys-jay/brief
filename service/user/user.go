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
func Register(user *model.User, isAdmin ...bool) error {
	// Hash password
	hash, salt, err := utility.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("could not hash user password, got error: %w", err)
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

	// ### Uncomment the lines below to retrieve token to be sent in email ###

	// token, err := middleware.CreateToken(user.ID, user.Email, user.Role)
	// if err != nil {
	// 	return fmt.Errorf("could not create token, got error: %w", err)
	// }

	// TODO: Send Verification Email to User
	// err = utility.SendMail(user.Email, "index.html", &utility.EmailData{Token: token, Subject: "Test"})
	// if err != nil {
	// 	return fmt.Errorf("could not send email, got error: %w", err)
	// }

	db := postgres.GetDB()
	err = db.CreateUser(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("could not create user, got error: %w", err)
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return nil
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

	token, err := middleware.CreateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("could not create token")
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	// Ensure that user is verified and not locked
	if !user.IsVerified || user.IsLocked {
		var err error
		if !user.IsVerified {
			err = fmt.Errorf("user not verified")
		} else {
			err = fmt.Errorf("cannot login, user is currently locked")
		}
		return nil, err
	}

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
	err := db.UpdateUser(context.TODO(), id, user)
	if err != nil {
		return fmt.Errorf("could not update user, got error: %w", err)
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return nil
}

// Verify contains business logic to verify a user
func Verify(id string) (*model.User, error) {
	db := postgres.GetDB()
	user, err := db.VerifyUser(context.TODO(), id)
	if err != nil {
		return nil, fmt.Errorf("could not verify user, got error: %w", err)
	}

	// Omit password and salt from response
	user.Password = ""
	user.Salt = ""

	return user, nil
}

// ResetPassword contains business logic to reset a user's password
func ResetPassword(id string, rp *model.ResetPassword) (*model.User, error) {
	hashedPassword, salt, err := utility.HashPassword(rp.Password)
	if err != nil {
		return nil, fmt.Errorf("could not hash password, got error: %w", err)
	}

	// Set hashed password and salt
	rp.Password = hashedPassword
	rp.Salt = salt

	db := postgres.GetDB()
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
