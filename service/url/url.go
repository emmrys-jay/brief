package url

import (
	"brief/internal/constant"
	"brief/internal/model"
	"brief/pkg/repository/storage/postgres"
	"brief/utility"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Redirect contains business logic to redirect a shortened url to the original url
func Redirect(hash string) (*model.URL, error) {

	db := postgres.GetDB()
	url, err := db.GetURL(context.TODO(), hash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("url not found")
		}
		return nil, fmt.Errorf("could not fetch url, got error %w", err)
	}

	return url, nil
}

// ADMIN & USER

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
			return fmt.Errorf("oops, '%s' already exists", url.Hash)
		}
	}

	return nil
}

// Delete contains business logic to delete a user's saved URL or a random url by its 'id'
func Delete(ctxInfo *model.ContextInfo, urlId string) (*model.URL, error) {

	db := postgres.GetDB()
	if ctxInfo.Role != constant.Roles[constant.Admin] {
		url, err := db.GetURLById(context.TODO(), urlId)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, fmt.Errorf("url not found")
			}
			return nil, fmt.Errorf("could not fetch url, got error %w", err)
		}

		if url.UserID != ctxInfo.ID {
			return nil, fmt.Errorf("unauthorized to perform this action")
		}
	}

	url, err := db.DeleteUrl(context.TODO(), urlId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("url not found")
		}
		return nil, fmt.Errorf("could not delete url, got error %w", err)
	}

	return url, nil
}

// GetURLs contains business logic to fetch all URL's created by a user with 'userID'
func GetURLs(userID string) ([]model.URL, error) {

	db := postgres.GetDB()
	urls, err := db.GetUrls(context.TODO(), userID)
	if err != nil {
		return nil, fmt.Errorf("could not get urls, got error : %w", err)
	}

	return urls, nil
}

// ADMIN

// GetAll contains business logic to fetch all URL's
func GetAll() ([]model.URL, error) {

	db := postgres.GetDB()
	urls, err := db.GetAll(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("could not get urls, got error : %w", err)
	}

	return urls, nil
}
