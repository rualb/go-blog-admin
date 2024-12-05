package service

import (
	"encoding/base64"
	"fmt"
	"go-blog-admin/internal/config"
	"time"
)

const (

// secretKeySize is len of secrets
// secretKeySize = 64
)

type VaultKey struct {
	ID        string `gorm:"size:255;primaryKey"` // `gorm:"size:255;primaryKey"`
	CreatedAt time.Time
	AuthKey   string `gorm:"size:255"` // base64 for auth
}

func (x VaultKey) IsEmpty() bool {
	return x.ID == "" || x.AuthKey == ""
}

type VaultService interface {
	CurrentKey() (secret *SecretKey, err error)
	KeyByID(id string) (secret *SecretKey, err error)

	KeyScopeAuth() VaultKeyScope

	// Append(secret ...SecretKey)
}

type defaultVaultService struct {
	keychain []SecretKey
}

func (x *defaultVaultService) KeyScopeAuth() VaultKeyScope {
	return &vaultKeyScope{
		vaultService: x,
		auth:         true,
	}
}

type VaultKeyScope interface {
	CurrentKey() (id string, secret []byte, err error)
	KeyByID(id string) (secret []byte, err error)
}

type vaultKeyScope struct {
	vaultService VaultService
	auth         bool
}

func (x vaultKeyScope) extractSecret(secret *SecretKey) ([]byte, error) {

	if x.auth {
		return secret.AuthKey, nil
	}

	return nil, fmt.Errorf("error set secret key type for extract")
}

func (x vaultKeyScope) CurrentKey() (id string, secret []byte, err error) {

	key, err := x.vaultService.CurrentKey()

	if err == nil && key == nil {
		err = fmt.Errorf("error no any secret key")
	}
	if err != nil {
		return "", nil, err
	}

	id = key.ID
	secret, err = x.extractSecret(key)

	if err != nil {
		return "", nil, err
	}

	return id, secret, nil
}

func (x vaultKeyScope) KeyByID(id string) (secret []byte, err error) {

	key, err := x.vaultService.KeyByID(id)

	if err == nil && key == nil {
		err = fmt.Errorf("error no any secret key")
	}
	if err != nil {
		return nil, err
	}

	secret, err = x.extractSecret(key)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

type SecretKey struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	AuthKey   []byte //   for auth
}

func (x SecretKey) IsEmpty() bool {
	return x.ID == "" || len(x.AuthKey) == 0
}
func allKeys(appService AppService) (keys []config.AppConfigVaultKey, err error) {

	keys = []config.AppConfigVaultKey{} // nil

	{

		// 1

		keysConfig := appService.Config().Vault.Keys

		keys = append(keys, keysConfig...) // nil is same as empty array for append

	}

	{
		// 2

		keysDB := make([]VaultKey, 1, 10)
		res := appService.Repository().Driver().Order("created_at desc").Limit(10).Find(&keysDB)
		if res.Error != nil {
			return nil, res.Error
		}

		for _, v := range keysDB {
			keys = append(keys, config.AppConfigVaultKey{
				ID:      v.ID,
				AuthKey: v.AuthKey,
			})
		}
	}

	return keys, err
}

func newVaultService(appService AppService) (VaultService, error) {

	res := &defaultVaultService{
		keychain: []SecretKey{},
	}

	keys, err := allKeys(appService)
	if err != nil {
		return nil, err
	}
	err = res.loadKeys(keys)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (x *defaultVaultService) loadKeys(keys []config.AppConfigVaultKey) (err error) {

	for _, itm := range keys {

		// is empty
		if itm.IsEmpty() {
			// may be error
			continue
		}

		k := SecretKey{}
		k.ID = itm.ID

		if k.AuthKey, err = base64.StdEncoding.DecodeString(itm.AuthKey); err != nil {
			return fmt.Errorf("error on un-base64 key %v :%v", itm.ID, err)
		}

		x.keychain = append(x.keychain, k)
	}

	return nil
}

func (x *defaultVaultService) CurrentKey() (secret *SecretKey, err error) {
	if len(x.keychain) > 0 {
		r := &x.keychain[len(x.keychain)-1]
		return r, nil
	}

	return nil, fmt.Errorf("error no any key")
}

func (x *defaultVaultService) KeyByID(id string) (secret *SecretKey, err error) {

	// TODO may be use map
	for _, itm := range x.keychain {
		if itm.ID == id {
			return &itm, nil
		}
	}

	return nil, fmt.Errorf("error key not exists: %v", id)
}
