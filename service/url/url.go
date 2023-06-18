package url

import (
	"brief/internal/model"
)

// Link contains business logic to shorten and store a URL
func Shorten(url *model.URL) error {
	if url.Hash == "" {

	}
	return nil
}

// GetURLs contains business logic to fetch all URL's stored by a user
func GetURLs(userID string) ([]model.URL, error) {

	return nil, nil
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
