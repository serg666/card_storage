package models

import (
	"time"

	vault "github.com/mittwald/vaultgo"
)

var (
	SESSION_LIFE_TIME_SECONDS = 20 * 60
)

type Session struct {
	ID       uint
	Key      string
	Body     string
	ExpireAt time.Time
	CKey     *Key
}

func (s Session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}

func (s *Session) EncryptedBody(vaultClient *vault.Client) (error, *string) {
	transit := vaultClient.Transit()

	encryptResponse, err := transit.Encrypt(s.CKey.Name, &vault.TransitEncryptOptions{
		Plaintext: s.Body,
	})

	if err != nil {
		return err, nil
	}

	return nil, &encryptResponse.Data.Ciphertext
}

func (s *Session) DecryptedBody(vaultClient *vault.Client) (error, *string) {
	transit := vaultClient.Transit()

	decryptResponse, err := transit.Decrypt(s.CKey.Name, &vault.TransitDecryptOptions{
		Ciphertext: s.Body,
	})

	if err != nil {
		return err, nil
	}

	return nil, &decryptResponse.Data.Plaintext
}

type GormSession struct {
	ID       uint      `gorm:"primarykey"`
	Key      string    `gorm:"type:varchar(255);unique_index"`
	Body     string    `gorm:"type:text"`
	ExpireAt time.Time `gorm:"index"`
	CKeyID   uint
	CKey     GormKey `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

func (sg GormSession) ToSession(vaultClient *vault.Client) (*Session, error) {
	session := &Session{
		ID:       sg.ID,
		Key:      sg.Key,
		Body:     sg.Body,
		ExpireAt: sg.ExpireAt,
	}
	if sg.CKey.ID > 0 {
		session.CKey = &Key{
			ID:   sg.CKey.ID,
			Name: sg.CKey.Name,
		}
		err, decryptedBody := session.DecryptedBody(vaultClient)

		if err != nil {
			return nil, err
		}

		session.Body = *decryptedBody
	}
	return session, nil
}

func SessionToGormSession(s *Session, vaultClient *vault.Client) (*GormSession, error) {
	body := s.Body

	if s.CKey != nil {
		err, encryptedBody := s.EncryptedBody(vaultClient)

		if err != nil {
			return nil, err
		}

		body = *encryptedBody
	}

	sg := &GormSession{
		Key:      s.Key,
		Body:     body,
		ExpireAt: s.ExpireAt,
	}

	if s.CKey != nil {
		sg.CKeyID = s.CKey.ID
		sg.CKey = GormKey{
			ID:   s.CKey.ID,
			Name: s.CKey.Name,
		}
	}

	return sg, nil
}

func NewSession(key, body string, ckey *Key) *Session {
	return &Session{
		Key:      key,
		Body:     body,
		ExpireAt: time.Now().Add(time.Second * time.Duration(SESSION_LIFE_TIME_SECONDS)),
		CKey:     ckey,
	}
}
