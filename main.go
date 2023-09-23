package main

import (
	vault "github.com/mittwald/vaultgo"
	"github.com/sirupsen/logrus"
	"github.com/serg666/card_storage/internal"
	"github.com/serg666/card_storage/models"
	"github.com/serg666/card_storage/pkg"
)

var (
	cfg  = pkg.NewConfigsByEnv()
	args = pkg.ReadArguments()
)

func encryptUnencrypted(cardRepo internal.CardRepository, sessionRepo internal.SessionRepository, keyRepo internal.KeyRepository) {
	if cards, err := cardRepo.Query(internal.NewCardSpecificationWithoutKey()); err == nil {
		for _, card := range cards {
			if randomKey, err := keyRepo.Query(internal.NewKeySpecificationWithRandom()); err == nil {
				if len(randomKey) == 0 {
					logrus.Warningf("No keys has been found")
					continue
				}
				card.Key = randomKey[0]
				if err := cardRepo.Update(card); err == nil {
					logrus.Infof("Success encrypt card (%d): %s", card.ID, card.PanMask())
				} else {
					logrus.Warningf("Error updating card (%d) %s: %v", card.ID, card.PanMask(), err)
				}
			} else {
				logrus.Warningf("Error quering random key: %v", err)
			}
		}
	} else {
		logrus.Warningf("Error quering card repo: %v", err)
	}

	if sessions, err := sessionRepo.Query(internal.NewSessionSpecificationWithoutKey()); err == nil {
		for _, session := range sessions {
			if randomKey, err := keyRepo.Query(internal.NewKeySpecificationWithRandom()); err == nil {
				if len(randomKey) == 0 {
					logrus.Warningf("No keys has been found")
					continue
				}
				session.CKey = randomKey[0]
				if err := sessionRepo.Update(session); err == nil {
					logrus.Infof("Success encrypt session: %d", session.ID)
				} else {
					logrus.Warningf("Error updating session (%d): %v", session.ID, err)
				}
			} else {
				logrus.Warningf("Error quering random key: %v", err)
			}
		}
	} else {
		logrus.Warningf("Error quering session repo: %v", err)
	}
}

func addKeys(keyRepo internal.KeyRepository, vaultClient *vault.Client) {
	transit := vaultClient.Transit()

	for i := 0; i < 1000; i++ {
		keyName, err := models.GenerateRandomBytes(32)
		if err != nil {
			logrus.Warningf("Error creating DEK name: %v", err)
			continue
		}
		key := &models.Key{
			Name: keyName,
		}
		if err := transit.Create(key.Name, &vault.TransitCreateOptions{
			Exportable:           vault.BoolPtr(false),
			AllowPlaintextBackup: vault.BoolPtr(false),
			Derived:              vault.BoolPtr(false),
			Type:                 "aes256-gcm96",
		}); err == nil {
			if m, err := keyRepo.Create(key); err == nil {
				logrus.Infof("Success creating %s DEK", m.Name)
			} else {
				logrus.Warningf("Error creating %s DEK: %v", key.Name, err)
			}
		} else {
			logrus.Warningf("Error creating %s DEK: %v", key.Name, err)
		}
	}
}

func main() {
	vaultClient, err := vault.NewClient(cfg.VaultAddress, vault.WithCaPath(cfg.VaultCapath), vault.WithAuthToken(cfg.VaultToken))
	if err != nil {
		logrus.Errorf("Failed to get Vault client: %v", err)
	}

	dbConn := pkg.NewDBConnector(cfg.DbDSN, cfg.SafeLogs)
	defer dbConn.Close()
	if err := dbConn.Open(); err != nil {
		logrus.Errorf("Failed on connect to DB: %v", err)
	}

	cardRepo := internal.NewGormCardRepository(dbConn.DB, vaultClient)
	sessionRepo := internal.NewGormSessionRepository(dbConn.DB, vaultClient)
	keyRepo := internal.NewGormKeyRepository(dbConn.DB)

	switch args.Run {
	case "add-keys":
		addKeys(keyRepo, vaultClient)
		logrus.Info("Add keys has been completed")
	case "encrypt-unencrypted":
		encryptUnencrypted(cardRepo, sessionRepo, keyRepo)
		logrus.Info("Encryption of unencrypded cards and sessions has been completed")
	case "migrate-up":
		logrus.Infof("Migrate up results: %v", dbConn.Migrate())
	case "migrate-down-all":
		logrus.Infof("Drop tables results: %v", dbConn.DropTables())
	case "delete-expired-sessions":
		logrus.Infof("Delete expired sessions results: %v", sessionRepo.DeleteExpired())
	case "app":
		if err := internal.NewServer(cfg, cardRepo, sessionRepo, keyRepo).Run(); err != nil {
			logrus.Errorf("Failed on run server: %v", err)
		}
	}
}
