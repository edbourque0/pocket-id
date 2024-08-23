package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserAuthorizedOidcClient struct {
	Scope  string
	UserID string `gorm:"primary_key;"`
	User   User

	ClientID string `gorm:"primary_key;"`
	Client   OidcClient
}

type OidcAuthorizationCode struct {
	Base

	Code      string
	Scope     string
	Nonce     string
	ExpiresAt time.Time

	UserID string
	User   User

	ClientID string
}

type OidcClient struct {
	Base

	Name         string
	Secret       string
	CallbackURLs CallbackURLs
	ImageType    *string
	HasLogo      bool `gorm:"-"`

	CreatedByID string
	CreatedBy   User
}

func (c *OidcClient) AfterFind(_ *gorm.DB) (err error) {
	// Compute HasLogo field
	c.HasLogo = c.ImageType != nil && *c.ImageType != ""
	return nil
}

type CallbackURLs []string

func (cu *CallbackURLs) Scan(value interface{}) error {
	if v, ok := value.([]byte); ok {
		return json.Unmarshal(v, cu)
	} else {
		return errors.New("type assertion to []byte failed")
	}
}

func (cu CallbackURLs) Value() (driver.Value, error) {
	return json.Marshal(cu)
}
