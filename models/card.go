package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	vault "github.com/mittwald/vaultgo"
)

var (
	CARD_TOKEN_LENGTH  = 32
	EXPIRE_DATE_FORMAT = "06/01"
)

type Card struct {
	ID      uint
	Token   string
	PAN     string
	Holder  *string
	ExpDate *time.Time
	Key     *Key
}

func GenerateRandomBytes(size int) (string, error) {
	h := sha256.New()
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	h.Write(b)
	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}

func (c *Card) GenerateToken() error {
	token, err := GenerateRandomBytes(CARD_TOKEN_LENGTH)
	if err != nil {
		return err
	}
	c.Token = token
	return nil
}

func (c *Card) EncryptedPan(vaultClient *vault.Client) (string, error) {
	transit := vaultClient.Transit()

	encryptResponse, err := transit.Encrypt(c.Key.Name, &vault.TransitEncryptOptions{
		Plaintext: c.PAN,
	})

	if err != nil {
		return "", err
	}

	return encryptResponse.Data.Ciphertext, err
}

func (c *Card) DecryptedPan(vaultClient *vault.Client) (string, error) {
	transit := vaultClient.Transit()

	decryptResponse, err := transit.Decrypt(c.Key.Name, &vault.TransitDecryptOptions{
		Ciphertext: c.PAN,
	})

	if err != nil {
		return "", err
	}

	return decryptResponse.Data.Plaintext, err
}

func (c Card) PanMask() string {
	return GetPanMask(c.PAN)
}

func (c Card) FormatedExpDate() *string {
	if c.ExpDate == nil {
		return nil
	}
	exp_formated := c.ExpDate.Format(EXPIRE_DATE_FORMAT)
	return &exp_formated
}

type GormCard struct {
	ID        uint   `gorm:"primarykey"`
	Token     string `gorm:"type:varchar(100);index:idx_token,unique"`
	PAN       string `gorm:"type:varchar(255);index"`
	MASK      string `gorm:"type:varchar(12);index"`
	KeyID     uint
	Key       GormKey    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Holder    *string    `gorm:"type:varchar(50);default:null"`
	ExpDate   *time.Time `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (cg GormCard) ToCard(vaultClient *vault.Client) (*Card, error) {
	card := &Card{
		ID:      cg.ID,
		Token:   cg.Token,
		PAN:     cg.PAN,
		Holder:  cg.Holder,
		ExpDate: cg.ExpDate,
	}
	if cg.Key.ID > 0 {
		card.Key = &Key{
			ID:   cg.Key.ID,
			Name: cg.Key.Name,
		}
		decryptedPan, err := card.DecryptedPan(vaultClient)

		if err != nil {
			return nil, err
		}

		card.PAN = decryptedPan
	}

	return card, nil
}

func CardToGormCard(c *Card, vaultClient *vault.Client) (*GormCard, error) {
	pan := c.PAN

	if c.Key != nil {
		encryptedPan, err := c.EncryptedPan(vaultClient)

		if err != nil {
			return nil, err
		}

		pan = encryptedPan
	}

	cg := &GormCard{
		Token:   c.Token,
		PAN:     pan,
		MASK:    c.PanMask(),
		Holder:  c.Holder,
		ExpDate: c.ExpDate,
	}

	if c.Key != nil {
		cg.KeyID = c.Key.ID
		cg.Key = GormKey{
			ID:   c.Key.ID,
			Name: c.Key.Name,
		}
	}

	return cg, nil
}

func CleanPan(pan string) string {
	re, _ := regexp.Compile(`[^\d]`)
	return re.ReplaceAllString(pan, "")
}

func GetPanMask(row string) string {
	return fmt.Sprintf("%s**%s", row[0:6], row[len(row)-4:])
}

func HidePan(row string) string {
	return regexp.MustCompile(`(\d{6})\d+(\d{4})`).ReplaceAllString(row, "$1******$2")
}
