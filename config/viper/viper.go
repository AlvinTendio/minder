package viper

import (
	"encoding/base64"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/AlvinTendio/minder/config"
	"github.com/AlvinTendio/minder/config/internal"
	viperpit "github.com/ajpauwels/pit-of-vipers"
	"github.com/spf13/viper"
)

type (
	Config struct {
		watchers       []internal.Watcher
		data           *viper.Viper
		additionalPath []string
		rootKey        string
		dataCh         chan *viper.Viper
		stopCh         chan struct{}

		// format defines the file format, whether it's JSON, YAML, etc. If the file's suffix is already specified,
		// there's no need to add this format, as Viper automatically detects it.
		format string
	}

	Option func(*Config)
)

// WithPrefix sets prefix key to get config from Viper.
// Viper supports accessing keys using dot notation. Add in a rootKey
// if you have a nested file with a root element. For e.g:
// { data: {"key1": "value1"}}, pass in "data" as rootKey to be able
// to be able to reference the inner layer keys
func WithPrefix(rootKey string) Option {
	return func(cfg *Config) {
		if rootKey != "" {
			cfg.rootKey += rootKey + "."
		}
	}
}

// WithAdditionalPath adds additional file path to the Viper list of files for obtaining configuration values and
// monitoring changes.
func WithAdditionalPath(paths []string) Option {
	return func(cfg *Config) {
		if len(paths) > 0 {
			cfg.additionalPath = paths
		}
	}
}

// WithKVVersion2 sets the prefix key for the default Version 2 Vault data structure.
func WithKVVersion2() Option {
	return WithPrefix("data")
}

// WithFormat will specify the file format to be added to Viper.
func WithFormat(format string) Option {
	return func(cfg *Config) {
		if format != "" {
			cfg.format = format
		}
	}
}

func NewConfig(pathFile string, options ...Option) *Config {
	// Declare new config and apply all options
	c := &Config{dataCh: make(chan *viper.Viper), stopCh: make(chan struct{})}
	for _, opt := range options {
		opt(c)
	}

	// Get all config path that specifies path to config file
	allPath := []string{pathFile}
	if len(c.additionalPath) > 0 {
		allPath = combineAllPath(allPath, c.additionalPath)
	}

	// Run the init() pass all the vipers instance, wait data for the first time
	// and run goroutine to check all possible update
	go c.init(allPath)
	c.data = <-c.dataCh
	go c.updateConfig()

	return c
}

// Internal function to combine all path
func combineAllPath(listPath ...[]string) []string {
	allPath := make([]string, 0)
	for _, list := range listPath {
		allPath = append(allPath, list...)
	}
	return allPath
}

func (c *Config) init(allPath []string) {
	// Make all possible viper instance. detect from allPath
	vipers := make([]*viper.Viper, len(allPath))
	for index, p := range allPath {
		v := c.initiateViper(p)

		vipers[index] = v
	}

	vpCh, errCh := viperpit.New(vipers)
	for {
		select {
		// Watch any possible changes from all viper instance
		// if any changes happen, pass AllSettings() to dataCh in order to communicate w/ other goroutine
		case vp := <-vpCh:
			log.Println("Config in watch", vp.AllKeys())
			c.dataCh <- vp

			// Log all possible error
		case err := <-errCh:
			log.Println(err)

			// Watch any stop signal from stop channel
		case <-c.stopCh:
			log.Println("init stopped!")
			return
		}
	}
}

func (c *Config) initiateViper(filePath string) *viper.Viper {
	v := viper.New()
	filename := path.Base(filePath)
	filePath = path.Dir(filePath)
	v.AddConfigPath(filePath)
	v.SetConfigName(filename)

	if c.format != "" {
		v.SetConfigType(c.format)
	}

	return v
}

func (c *Config) updateConfig() {
	for {
		select {
		// Wait any update from data channel, if update happen then renew data fields in memory
		case data := <-c.dataCh:
			c.data = data
			log.Println("Config updated :", c.data.AllKeys())

			// Wait any signal from stop channel, if signal occur then stop the goroutine
		case <-c.stopCh:
			log.Println("update config stopped!")
			return
		}
	}
}

// Get returns configuration value as interface
func (c *Config) Get(key string) interface{} {
	return c.data.Get(c.rootKey + key)
}

// GetInt returns configuration value as integer 64 bit
func (c *Config) GetInt(key string) int64 {
	return c.data.GetInt64(c.rootKey + key)
}

// GetString returns configuration value as string
func (c *Config) GetString(key string) string {
	return c.data.GetString(c.rootKey + key)
}

// GetBool returns configuration value as boolean
func (c *Config) GetBool(key string) bool {
	return c.data.GetBool(c.rootKey + key)
}

// GetFloat returns configuration value as float 64 bit
func (c *Config) GetFloat(key string) float64 {
	return c.data.GetFloat64(c.rootKey + key)
}

// GetBinary returns configuration value as byte array,
// configuration value is stored as base64 encoded
func (c *Config) GetBinary(key string) []byte {
	value := c.data.Get(c.rootKey + key)
	if value != nil {
		strValue, ok := value.(string)
		if !ok {
			log.Println("unable to convert to string")
			return nil
		}
		bytes, err := base64.StdEncoding.DecodeString(strValue)
		if err == nil {
			return bytes
		}
	}
	return nil
}

// GetArray returns configuration value as array
// configuration value is stored with format <element1>,<element2>,...
func (c *Config) GetArray(key string) []string {
	value := c.data.GetString(c.rootKey + key)

	if value == "" {
		return nil
	}

	return strings.Split(value, ",")
}

// GetMap returns configuration value as map
// configuration value is stored with format <key1>:<value1>,<key2>:<value2>,...
func (c *Config) GetMap(key string) map[string]string {
	str := c.GetString(c.rootKey + key)
	arr := strings.Split(str, ",")

	maps := make(map[string]string)

	kvLen := 2
	for _, element := range arr {
		kv := strings.SplitN(element, ":", 2)
		if len(kv) == kvLen {
			maps[kv[0]] = kv[1]
		}
	}
	return maps
}

func (c *Config) Close() error {
	for _, watcher := range c.watchers {
		watcher.Close()
	}
	close(c.stopCh)
	return nil
}

func (c *Config) Watch(keys ...string) <-chan []string {
	watcher := internal.NewWatcher(keys, c)
	c.watchers = append(c.watchers, watcher)
	return watcher.Change()
}

func (c *Config) GetCredentials(key string) config.Credentials {
	user := c.GetString(fmt.Sprintf("%s.user", key))
	password := c.GetString(fmt.Sprintf("%s.password", key))

	return config.Credentials{
		User:     user,
		Password: password,
	}
}

func (c *Config) GetTLSCertificate(key string) config.TLSCertificate {
	ca := c.GetBinary(fmt.Sprintf("%s.ca", key))
	cert := c.GetBinary(fmt.Sprintf("%s.cert", key))
	pkey := c.GetBinary(fmt.Sprintf("%s.key", key))

	return config.TLSCertificate{
		CA:          ca,
		Certificate: cert,
		Key:         pkey,
	}
}
