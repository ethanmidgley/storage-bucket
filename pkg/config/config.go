package config

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

// Conf global variable for storing the config from the file
var Conf *Config

// Config global config variable
type Config struct {
	Yaml    *YamlConfig
	KeysMap map[string]bool
}

// YamlConfig stuct for yaml file
type YamlConfig struct {
	ControlPlane struct {
		Host           string   `yaml:"host"`
		Port           string   `yaml:"port"`
		NumberOfKeys   int      `yaml:"number-of-keys"`
		Keys           []string `yaml:"keys"`
		AllowedOrigins []string `yaml:"allowed-origins"`
	} `yaml:"control-plane"`
	Bucket struct {
		Location string `yaml:"filename"`
		Export   struct {
			Allowed     bool   `yaml:"allowed"`
			Compression string `yaml:"compression"`
		} `yaml:"export"`
	} `yaml:"bucket"`
}

// Load config from the config.yaml file
func Load() (*Config, error) {
	c := &Config{KeysMap: make(map[string]bool)}
	y := &YamlConfig{}

	file, err := os.Open("bucket.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&y); err != nil {
		return nil, err
	}

	c.Yaml = y

	log.Println("Checking api keys")
	c.CheckKey()

	// load the keys to from the array to a map
	for _, key := range c.Yaml.ControlPlane.Keys {
		c.KeysMap[key] = true
	}

	return c, nil
}

// CheckKey will make sure there is a key available
func (c *Config) CheckKey() {
	keysMissing := c.Yaml.ControlPlane.NumberOfKeys - len(c.Yaml.ControlPlane.Keys)
	if keysMissing > 0 {
		log.Println("Not enough api keys found generating the rest")

		var keys []string

		for i := 0; i < keysMissing; i++ {
			key := GenerateNewKey()
			keys = append(keys, key)
			h := sha256.New()
			h.Write([]byte(key))
			c.Yaml.ControlPlane.Keys = append(c.Yaml.ControlPlane.Keys, hex.EncodeToString(h.Sum(nil)))
		}

		// we are going to save the key to a tempary file so we can use it
		f, err := os.Create("api-keys.txt")
		if err != nil {
			log.Panic(err)
		}
		defer f.Close()

		_, err = f.Write([]byte("Please delete this file once the key has been copied\nKeys:\n" + strings.Join(keys, "\n")))
		if err != nil {
			log.Panic(err)
		}

		c.Update()
		log.Println("Missing keys generated")
	} else {
		log.Println("All keys found")
	}

}

// Update will take the current stuct and turn it back to the yaml file
func (c *Config) Update() {

	y := c.Yaml

	ce, err := yaml.Marshal(&y)
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile("bucket.yaml", ce, 0644)

}

// GenerateNewKey will create a api key
func GenerateNewKey() string {

	u := uuid.New()
	return u.String()

}
