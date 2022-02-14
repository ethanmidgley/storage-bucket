package config

import (
	"io/ioutil"
	"log"
	"os"

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
		HTPasswd       string   `yaml:"htpasswd"`
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
	if c.CheckKey() {
		// load the keys to from the array to a map
		for _, key := range c.Yaml.ControlPlane.Keys {
			c.KeysMap[key] = true
		}
		log.Println("Keys found")
	} else {
		log.Printf("Keys were not found send a post request http://%s:%s/generate with your username and password set in the config to create api keys", c.Yaml.ControlPlane.Host, c.Yaml.ControlPlane.Port)
	}

	return c, nil
}

// CheckKey will make sure there is a key available
func (c *Config) CheckKey() bool {
	keysMissing := c.Yaml.ControlPlane.NumberOfKeys - len(c.Yaml.ControlPlane.Keys)
	if keysMissing > 0 {
		return false
	}
	return true

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
