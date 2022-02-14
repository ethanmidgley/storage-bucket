package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/ethanmidgley/storage-bucket/pkg/config"
	"github.com/google/uuid"
)

// GenerateKeys will generate uuid v4 keys based of the config yaml file
func GenerateKeys() ([]string, []string) {

	var keys []string
	var keyshashed []string

	for i := 0; i < config.Conf.Yaml.ControlPlane.NumberOfKeys; i++ {
		key := GenerateNewKey()
		keys = append(keys, key)
		h := sha256.New()
		h.Write([]byte(key))
		keyshashed = append(keyshashed, hex.EncodeToString(h.Sum(nil)))
	}

	return keys, keyshashed

}

// GenerateNewKey will create a api key
func GenerateNewKey() string {

	u := uuid.New()
	return u.String()

}
