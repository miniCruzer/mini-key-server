package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type IKey interface {
	ToHex() string
	WriteSecret(size int) (int, error)
}

type Key struct {
	gorm.Model
	ApplicationID uint
	Enabled       bool
	Memo          *string
	HardwareID    *string
	Remaining     uint
	Secret        string
}

func (k *Key) SetRandomSecret(size int) (n int, err error) {
	b := make([]byte, size)
	n, err = rand.Read(b)

	if err != nil {
		return
	}

	k.Secret = base64.StdEncoding.EncodeToString(b)
	return
}

func (k *Key) GetSecretBytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(k.Secret)
}

func (k *Key) FromJWT(myToken string, keyLookup jwt.Keyfunc) error {
	token, err := jwt.ParseWithClaims(myToken, &KeyClaims{}, keyLookup)

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("JWT not valid")
	}

	claims, _ := token.Claims.(*KeyClaims)
	k.ID = claims.ID

	return nil
}

func lookupKey(db *gorm.DB, keyID uint) ([]byte, error) {
	var key Key
	err := db.First(&key, keyID).Error

	if err != nil {
		return nil, err
	}

	return key.GetSecretBytes()
}

type KeyClaims struct {
	jwt.StandardClaims
	ID uint `json:"id"`
}

type KeyActivation struct {
	gorm.Model
	Key        Key
	Identifier *string
}