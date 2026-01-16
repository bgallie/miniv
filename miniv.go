// Copyright 2025 Billy G. Allie

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package miniv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
)

// Miniv represents the configuration management structure.
// It provides methods to set and get configuration values,
// bind flags, manage environment variables, and handle defaults.
// It also supports reading from and writing to a TOML configuration file.
type Miniv struct {
	configPath          string
	configFile          string
	envPrefix           string
	automaticEnvApplied bool
	emptyEnvVarValid    bool
	setvalues           map[string]any
	boundFlags          map[string]any
	envVars             map[string]string
	cfgValues           map[string]any
	flatCfgValues       map[string]any
	defaults            map[string]any
}

// NewConfig creates a new Config instance with default values.
// ConfigPath is set to the current working directory.
// ConfigFile is set to "config.toml".
// EnvPrefix is set to an empty string.
// The internal maps are initialized as empty maps.
func NewConfig() *Miniv {
	return &Miniv{
		automaticEnvApplied: false,
		setvalues:           make(map[string]any),
		boundFlags:          make(map[string]any),
		envVars:             make(map[string]string),
		cfgValues:           make(map[string]any),
		flatCfgValues:       make(map[string]any),
		defaults:            make(map[string]any),
	}
}

// SetConfigPath sets the configuration file path.
// Default is current working directory.
func (c *Miniv) SetConfigPath(configPath string) {
	c.configPath = configPath
}

// GetConfigPath returns the configuration file path.
// Default is current working directory.
func (c *Miniv) GetConfigPath() string {
	return c.configPath
}

// SetConfigFile sets the configuration file name.
// Default is "config.toml".
func (c *Miniv) SetConfigFile(configFile string) {
	c.configFile = configFile
}

// GetConfigFile returns the configuration file name.
// Default is "config.toml".
func (c *Miniv) GetConfigFile() string {
	return c.configFile
}

// SetEnvPrefix sets the environment variable prefix.
// Default is empty string.
func (c *Miniv) SetEnvPrefix(envPrefix string) {
	c.envPrefix = envPrefix
}

// GetEnvPrefix returns the environment variable prefix.
// Default is empty string.
func (c *Miniv) GetEnvPrefix() string {
	return c.envPrefix
}

// SetEmptyEnvVarValid sets whether empty environment variables are considered
// valid.
func (c *Miniv) SetEmptyEnvVarValid(valid bool) {
	c.emptyEnvVarValid = valid
}

// GetEmptyEnvVarValid returns whether empty environment variables are
// considered valid.
func (c *Miniv) GetEmptyEnvVarValid() bool {
	return c.emptyEnvVarValid
}

// SetValue sets a value for a key.
// If the key already exists, it overwrites the existing value.
func (c *Miniv) SetValue(key string, value any) {
	c.setvalues[key] = value
}

// GetValue returns the set value for a key.
// If the key does not exist, returns nil and false.
func (c *Miniv) GetValue(key string) (any, bool) {
	val, exists := c.setvalues[key]
	return val, exists
}

// BindFlag binds a pflag.Flag to a config key.
// If the key already exists, it overwrites the existing flag.
// Example usage:
//
//	var exampleFlag = flag.String("example-flag", "default", "An example flag")
//	config.BindFlag("example.flag", flag.Lookup("example-flag"))
//
// This will bind the "example-flag" to the config.
// Subsequent calls to config.Get("example.flag") will return the value of the flag.
// Note: This does not automatically update the config values when the flag is changed.
// You need to call config.Get to retrieve the current value.
// Example:
//
//	flag.Parse()
//	value, _ := config.Get("example.flag")
//	fmt.Println("Example Flag Value:", value)
//
// Note: The flag must be defined in the pflag package.
func (c *Miniv) BindFlag(key string, flag *pflag.Flag) {
	c.boundFlags[key] = flag
}

// BindFlags binds all flags from a pflag.FlagSet to the config.
// It iterates over all flags in the FlagSet and binds each one.
// If a flag already exists for a key, it overwrites the existing flag.
// Example usage:
//
//	var flagSet = pflag.NewFlagSet("example", pflag.ExitOnError)
//	flagSet.String("example-flag", "default", "An example flag")
//	config.BindFlags(flagSet)
//
// This will bind the "example-flag" to the config.
// Subsequent calls to config.Get("example-flag") will return the value of the flag.
// Note: This does not automatically update the config values when flags are changed.
// You need to call config.Get to retrieve the current value.
// Example:
//
//	flagSet.Parse(os.Args[1:])
//	value, _ := config.Get("example-flag")
//	fmt.Println("Example Flag Value:", value)
func (c *Miniv) BindFlags(flagSet *pflag.FlagSet) {
	flagSet.VisitAll(func(flag *pflag.Flag) {
		c.BindFlag(flag.Name, flag)
	})
}

// GetBoundFlag returns the bound flag for a key.
// If the key does not exist, returns nil and false.
func (c *Miniv) GetBoundFlag(key string) (*pflag.Flag, bool) {
	val, exists := c.boundFlags[key]
	if !exists {
		return nil, false
	}
	flag, ok := val.(*pflag.Flag)
	return flag, ok
}

// GetBoundFlagValue returns the value of a bound flag for a key.
// If the key does not exist or the flag is not changed, returns nil and false.
func (c *Miniv) GetBoundFlagValue(key string) (any, bool) {
	flag, exists := c.GetBoundFlag(key)
	if exists && flag.Changed {
		return flag.Value, true
	}
	return nil, false
}

// AutomaticEnv enables automatic environment variable binding.
// It will check for an environment variable any time a config.Get request is made.
// The environment variable name is created by uppercasing the key and
// replacing periods (.) with underscores (_).  For example, the key
// "database.url" would bind to the environment variable "DATABASE_URL".
func (v *Miniv) AutomaticEnv() {
	v.automaticEnvApplied = true
}

// SetEnvVar sets the environment variable name for a key.
// If the key already exists, it overwrites the existing env var name.
func (c *Miniv) SetEnvVar(key string, envVar string) {
	var transformedKey string
	if len(c.envPrefix) > 0 {
		transformedKey = strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%s_%s", c.envPrefix, key), ".", "_"))
	} else {
		transformedKey = strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	}
	c.envVars[transformedKey] = envVar
}

// GetEnvVar returns the environment variable value for a key.
// If the key does not exist, returns empty string and false.
// If AutomaticEnv is enabled, it first checks for the automatically generated
// envvar name.
// If not found, it checks for an explicitly set env var name in the envVars map.
func (c *Miniv) GetEnvVar(key string) (val string, exists bool) {
	if c.automaticEnvApplied {
		// Check for automatically generated env var name.
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		if len(c.envPrefix) > 0 {
			envKey = c.envPrefix + "_" + envKey
		}
		val, exists = os.LookupEnv(envKey)
		if exists {
			if !c.emptyEnvVarValid && val == "" {
				return "", false
			}
			return
		}
	}
	// Check for explicitly set env var name in the envVars map.
	var transformedKey string
	if len(c.envPrefix) > 0 {
		transformedKey = strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%s_%s", c.envPrefix, key), ".", "_"))
	} else {
		transformedKey = strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
	}
	if envVarName, exists := c.envVars[transformedKey]; exists {
		val, exists = os.LookupEnv(envVarName)
		if exists && !c.emptyEnvVarValid && val == "" {
			return "", false
		}
		return val, exists
	}
	return "", false
}

// SetDefault sets the default value for a key.
// If the key already exists, it overwrites the existing default.
func (c *Miniv) SetDefault(key string, value any) {
	c.defaults[key] = value
}

// GetDefault returns the default value for a key.
// If the key does not exist, returns nil and false.
func (c *Miniv) GetDefault(key string) (any, bool) {
	val, exists := c.defaults[key]
	return val, exists
}

// GetConfigValue returns a config value by key.
// The key should be in dot notation for nested values.
// If the key exists in cfgValues, returns that value.
// Otherwise, checks flatCfgValues for the key.
// If the key does not exist, returns nil and false.
func (c *Miniv) GetConfigValue(key string) (any, bool) {
	if val, exists := c.cfgValues[key]; exists {
		return val, exists
	}
	val, exists := c.flatCfgValues[key]
	if exists {
		return val, exists
	}
	return nil, false
}

// flattenConfigValues flattens nested configuration values into a flat map.
// It uses dot notation for nested keys.  If a top level key already contains
// a dot, it takes precedence and is not flattened. for example, given the
// following nested map:
//
//		{
//		  "database": {
//		    "url": "localhost",
//		    "port": 5432
//		  }
//		   "some.other.key": "value"
//	    "database.url": "remotehost"
//		}
//
// The resulting flat map would be:
//
//		{
//		  "database.url": "remotehost" #	"database.url": "localhost" is ignored
//		  "database.port": 5432
//	   "some.other.key": "value"
//		}
func (c *Miniv) flattenConfigValues(prefix string, values map[string]any, flatMap map[string]any) {
	for k, v := range values {
		if prefix != "" {
			// Avoid double dot notation if key already contains a dot.
			if strings.Contains(k, ".") {
				continue
			}
			k = prefix + "." + k
		} else {
			// Add keys that already contain dots only if no prefix,
			// overwriting existing values.
			if strings.Contains(k, ".") {
				flatMap[k] = v
				continue
			}
		}
		switch v := v.(type) {
		case map[string]any:
			c.flattenConfigValues(k, v, flatMap)
		case []any:
			for i, item := range v {
				itemKey := fmt.Sprintf("%s.%d", k, i)
				switch item := item.(type) {
				case map[string]any:
					c.flattenConfigValues(itemKey, item, flatMap)
				default:
					flatMap[itemKey] = item
				}
			}
		default:
			// Only set the value if the key does not already exist.
			// This prevents overwriting dot notation keys that were set
			// at the top level of the config.
			_, exists := flatMap[k]
			if !exists {
				flatMap[k] = v
			}
		}
	}
}

// Get returns the value for a key.
// The order of precedence is:
// 1. Set values via SetValue
// 2. Bound flag values via GetBoundFlagValue
// 3. Environment variable values via GetEnvVar
// 4. Flattened config file values via flatCfgValues
// 5. Default values via GetDefault
// If the key does not exist in any of these, returns nil and false.
func (c *Miniv) Get(key string) (any, bool) {
	if val, exists := c.setvalues[key]; exists {
		return val, exists
	}
	if val, exists := c.GetBoundFlagValue(key); exists {
		return val, exists
	}
	if val, exists := c.GetEnvVar(key); exists {
		return val, exists
	}
	if val, exists := c.flatCfgValues[key]; exists {
		return val, exists
	}
	if val, exists := c.defaults[key]; exists {
		return val, exists
	}
	return nil, false
}

// GetString returns the string value for a key.
// If the key does not exist or is not a string, returns an empty string.
func (c *Miniv) GetString(key string) string {
	val, _ := c.Get(key)
	return cast.ToString(val)
}

// GetStringSlice returns the string slice value for a key.
// If the key does not exist or is not a string slice, returns an empty slice.
func (c *Miniv) GetStringSlice(key string) []string {
	val, _ := c.Get(key)
	return cast.ToStringSlice(val)
}

// GetInt returns the int value for a key.
// If the key does not exist or is not an int, returns 0.
func (c *Miniv) GetInt64(key string) int64 {
	val, _ := c.Get(key)
	return cast.ToInt64(val)
}

// GetInt64Slice returns the int64 slice value for a key.
// If the key does not exist or is not an int64 slice, returns an empty slice.
func (c *Miniv) GetInt64Slice(key string) []int64 {
	val, _ := c.Get(key)
	return cast.ToInt64Slice(val)
}

// GetFloat64 returns the float64 value for a key.
// If the key does not exist or is not a float64, returns 0.0.
func (c *Miniv) GetFloat64(key string) float64 {
	val, _ := c.Get(key)
	return cast.ToFloat64(val)
}

// GetFloat64Slice returns the float64 slice value for a key.
// If the key does not exist or is not a float64 slice, returns an empty slice.
func (c *Miniv) GetFloat64Slice(key string) []float64 {
	val, _ := c.Get(key)
	return cast.ToFloat64Slice(val)
}

// GetBool returns the boolean value for a key.
// If the key does not exist or is not a boolean, returns false.
func (c *Miniv) GetBool(key string) bool {
	val, _ := c.Get(key)
	return cast.ToBool(val)
}

// GetBoolSlice returns the boolean slice value for a key.
// If the key does not exist or is not a boolean slice, returns an empty slice.
func (c *Miniv) GetBoolSlice(key string) []bool {
	val, _ := c.Get(key)
	return cast.ToBoolSlice(val)
}

// ConfigFileUsed returns the full path to the configuration file used.
// It combines configPath and configFile.
func (c *Miniv) ConfigFileUsed() string {
	return filepath.Clean(filepath.Join(c.configPath, c.configFile))
}

// ReadInConfig reads the configuration from the TOML file
// specified by configPath and configFile.
// It populates the cfgValues and flatCfgValues maps.
func (c *Miniv) ReadInConfig() error {
	if filepath.Ext(c.configFile) != ".toml" {
		return fmt.Errorf("invalid config file extension")
	}
	cfgIn, err := os.Open(c.ConfigFileUsed())
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer cfgIn.Close()
	d := toml.NewDecoder(cfgIn)
	if err = d.Decode(&c.cfgValues); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}
	c.flattenConfigValues("", c.cfgValues, c.flatCfgValues)
	return nil
}

// WriteConfig writes the configuration to the TOML file specified by
// configPath and configFile.
// It overwrites any existing file.
func (c *Miniv) WriteConfig() error {
	cfgFile := c.ConfigFileUsed()
	cfgOut, err := os.Create(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to create the config file: %w", err)
	}
	defer cfgOut.Close()
	enc := toml.NewEncoder(cfgOut)
	return enc.Encode(&c.cfgValues)
}

// WriteConfigAs writes the configuration to the specified TOML file.
// It overwrites any existing file.
func (c *Miniv) WriteConfigAs(cfgFile string) error {
	cfgOut, err := os.Create(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to create the config file: %w", err)
	}
	defer cfgOut.Close()
	enc := toml.NewEncoder(cfgOut)
	return enc.Encode(&c.cfgValues)
}

// SafeWriteConfig writes the configuration to the TOML file specified by
// configPath and configFile only if the file does not already exist.
// It prevents overwriting an existing configuration file.
func (c *Miniv) SafeWriteConfig() error {
	cfgFile := c.ConfigFileUsed()
	_, err := os.Stat(cfgFile)
	if err == nil {
		return fmt.Errorf("%s already exists.  Not overwriting", cfgFile)
	}
	cfgOut, err := os.Create(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to create the config file: %w", err)
	}
	defer cfgOut.Close()
	enc := toml.NewEncoder(cfgOut)
	return enc.Encode(&c.cfgValues)
}

// SafeWriteConfigAs writes the configuration to the specified TOML file only
// if the file does not already exist.
// It prevents overwriting an existing configuration file.
func (c *Miniv) SafeWriteConfigAs(cfgFile string) error {
	_, err := os.Stat(cfgFile)
	if err == nil {
		return fmt.Errorf("%s already exists.  Not overwriting", cfgFile)
	}
	cfgOut, err := os.Create(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to create the config file: %w", err)
	}
	defer cfgOut.Close()
	enc := toml.NewEncoder(cfgOut)
	return enc.Encode(&c.cfgValues)
}
