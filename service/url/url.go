package url

import (
	"brief/internal/model"
	"brief/pkg/repository/storage/postgres"
	"brief/utility"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Link contains business logic to shorten and store a URL
func Shorten(url *model.URL) error {
	url.ID = uuid.NewString()
	db := postgres.GetDB()

	if url.Hash == "" {
		for {
			hash, err := utility.GetURLHash(url.ID, url.LongURL)
			if err != nil {
				return fmt.Errorf("could not generate hash, got error %w", err)
			}
			url.Hash = hash

			if err := db.CreateURL(context.TODO(), url); err != nil {
				if !errors.Is(err, gorm.ErrDuplicatedKey) {
					return fmt.Errorf("could not store url, got error %w", err)
				}
			} else {
				break
			}
		}
	} else {
		if err := db.CreateURL(context.TODO(), url); err != nil {
			if !errors.Is(err, gorm.ErrDuplicatedKey) {
				return fmt.Errorf("could not store url, got error %w", err)
			}
			return fmt.Errorf("hash '%s' already exists", url.Hash)
		}
	}

	return nil
}

// GetURLs contains business logic to fetch all URL's stored by a user
func GetURLs(userID string) ([]model.URL, error) {

	db := postgres.GetDB()
	urls, err := db.GetUrls(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("could not get urls, got error : %w", err)
	}

	return urls, nil
}

// Delete contains business logic to delete a user's saved URL
func Delete(userId, urlId string) (*model.URL, error) {

	return nil, nil
}

// GetAll contains business logic to fetch all URL's
func GetAll() ([]model.URL, error) {

	return nil, nil
}

// DeleteUrlByID contains business logic to delete a random url sepecified by 'id'
func DeleteUrlByID(id string) (*model.URL, error) {

	return nil, nil
}
