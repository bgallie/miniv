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
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Miniv
	}{
		{
			name: "creates new config with defaults",
			want: &Miniv{
				automaticEnvApplied: false,
				setvalues:           make(map[string]any),
				boundFlags:          make(map[string]any),
				envVars:             make(map[string]string),
				cfgValues:           make(map[string]any),
				flatCfgValues:       make(map[string]any),
				defaults:            make(map[string]any),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig()
			if got.automaticEnvApplied != tt.want.automaticEnvApplied {
				t.Errorf("NewConfig().automaticEnvApplied = %v, want %v", got.automaticEnvApplied, tt.want.automaticEnvApplied)
			}
			if len(got.setvalues) != 0 {
				t.Errorf("NewConfig().setvalues should be empty, got %v", got.setvalues)
			}
			if len(got.boundFlags) != 0 {
				t.Errorf("NewConfig().boundFlags should be empty, got %v", got.boundFlags)
			}
			if len(got.envVars) != 0 {
				t.Errorf("NewConfig().envVars should be empty, got %v", got.envVars)
			}
			if len(got.cfgValues) != 0 {
				t.Errorf("NewConfig().cfgValues should be empty, got %v", got.cfgValues)
			}
			if len(got.flatCfgValues) != 0 {
				t.Errorf("NewConfig().flatCfgValues should be empty, got %v", got.flatCfgValues)
			}
			if len(got.defaults) != 0 {
				t.Errorf("NewConfig().defaults should be empty, got %v", got.defaults)
			}
		})
	}
}

func TestMiniv_SetConfigPath(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets config path",
			c:    NewConfig(),
			args: args{configPath: "/tmp"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetConfigPath(tt.args.configPath)
			if tt.c.configPath != tt.args.configPath {
				t.Errorf("SetConfigPath() did not set configPath correctly")
			}
		})
	}
}

func TestMiniv_GetConfigPath(t *testing.T) {
	tests := []struct {
		name string
		c    *Miniv
		want string
	}{
		{
			name: "gets default config path",
			c:    NewConfig(),
			want: "",
		},
		{
			name: "gets set config path",
			c: func() *Miniv {
				c := NewConfig()
				c.SetConfigPath("/tmp")
				return c
			}(),
			want: "/tmp",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetConfigPath(); got != tt.want {
				t.Errorf("Miniv.GetConfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_SetConfigFile(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets config file",
			c:    NewConfig(),
			args: args{configFile: "test.toml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetConfigFile(tt.args.configFile)
			if tt.c.configFile != tt.args.configFile {
				t.Errorf("SetConfigFile() did not set configFile correctly")
			}
		})
	}
}

func TestMiniv_GetConfigFile(t *testing.T) {
	tests := []struct {
		name string
		c    *Miniv
		want string
	}{
		{
			name: "gets default config file",
			c:    NewConfig(),
			want: "",
		},
		{
			name: "gets set config file",
			c: func() *Miniv {
				c := NewConfig()
				c.SetConfigFile("test.toml")
				return c
			}(),
			want: "test.toml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetConfigFile(); got != tt.want {
				t.Errorf("Miniv.GetConfigFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_SetEnvPrefix(t *testing.T) {
	type args struct {
		envPrefix string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets env prefix",
			c:    NewConfig(),
			args: args{envPrefix: "TEST"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetEnvPrefix(tt.args.envPrefix)
			if tt.c.envPrefix != tt.args.envPrefix {
				t.Errorf("SetEnvPrefix() did not set envPrefix correctly")
			}
		})
	}
}

func TestMiniv_GetEnvPrefix(t *testing.T) {
	tests := []struct {
		name string
		c    *Miniv
		want string
	}{
		{
			name: "gets default env prefix",
			c:    NewConfig(),
			want: "",
		},
		{
			name: "gets set env prefix",
			c: func() *Miniv {
				c := NewConfig()
				c.SetEnvPrefix("TEST")
				return c
			}(),
			want: "TEST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetEnvPrefix(); got != tt.want {
				t.Errorf("Miniv.GetEnvPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_SetEmptyEnvVarValid(t *testing.T) {
	type args struct {
		valid bool
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets empty env var valid to true",
			c:    NewConfig(),
			args: args{valid: true},
		},
		{
			name: "sets empty env var valid to false",
			c:    NewConfig(),
			args: args{valid: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetEmptyEnvVarValid(tt.args.valid)
			if tt.c.emptyEnvVarValid != tt.args.valid {
				t.Errorf("SetEmptyEnvVarValid() did not set emptyEnvVarValid correctly")
			}
		})
	}
}

func TestMiniv_GetEmptyEnvVarValid(t *testing.T) {
	tests := []struct {
		name string
		c    *Miniv
		want bool
	}{
		{
			name: "gets default empty env var valid",
			c:    NewConfig(),
			want: false,
		},
		{
			name: "gets set empty env var valid true",
			c: func() *Miniv {
				c := NewConfig()
				c.SetEmptyEnvVarValid(true)
				return c
			}(),
			want: true,
		},
		{
			name: "gets set empty env var valid false",
			c: func() *Miniv {
				c := NewConfig()
				c.SetEmptyEnvVarValid(false)
				return c
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetEmptyEnvVarValid(); got != tt.want {
				t.Errorf("Miniv.GetEmptyEnvVarValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_SetValue(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets string value",
			c:    NewConfig(),
			args: args{key: "test", value: "value"},
		},
		{
			name: "sets int value",
			c:    NewConfig(),
			args: args{key: "number", value: 42},
		},
		{
			name: "overwrites existing value",
			c: func() *Miniv {
				c := NewConfig()
				c.SetValue("test", "old")
				return c
			}(),
			args: args{key: "test", value: "new"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetValue(tt.args.key, tt.args.value)
			if val, exists := tt.c.setvalues[tt.args.key]; !exists || val != tt.args.value {
				t.Errorf("SetValue() did not set value correctly")
			}
		})
	}
}

func TestMiniv_GetValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *Miniv
		args  args
		want  any
		want1 bool
	}{
		{
			name: "gets existing value",
			c: func() *Miniv {
				c := NewConfig()
				c.SetValue("test", "value")
				return c
			}(),
			args:  args{key: "test"},
			want:  "value",
			want1: true,
		},
		{
			name:  "gets non-existing value",
			c:     NewConfig(),
			args:  args{key: "nonexistent"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.GetValue(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Miniv.GetValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMiniv_BindFlag(t *testing.T) {
	type args struct {
		key  string
		flag *pflag.Flag
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "binds a flag",
			c:    NewConfig(),
			args: args{
				key: "test-flag",
				flag: func() *pflag.Flag {
					fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
					fs.String("test-flag", "default", "test flag")
					return fs.Lookup("test-flag")
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.BindFlag(tt.args.key, tt.args.flag)
			if _, exists := tt.c.boundFlags[tt.args.key]; !exists {
				t.Errorf("BindFlag() did not bind flag correctly")
			}
		})
	}
}

func TestMiniv_BindFlags(t *testing.T) {
	type args struct {
		flagSet *pflag.FlagSet
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "binds multiple flags",
			c:    NewConfig(),
			args: args{
				flagSet: func() *pflag.FlagSet {
					fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
					fs.String("flag1", "default1", "flag1")
					fs.String("flag2", "default2", "flag2")
					return fs
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.BindFlags(tt.args.flagSet)
			if len(tt.c.boundFlags) != 2 {
				t.Errorf("BindFlags() did not bind all flags, got %d", len(tt.c.boundFlags))
			}
		})
	}
}

func TestMiniv_GetBoundFlag(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *Miniv
		args  args
		want  *pflag.Flag
		want1 bool
	}{
		{
			name: "gets bound flag",
			c: func() *Miniv {
				c := NewConfig()
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("test", "default", "test flag")
				flag := fs.Lookup("test")
				flag.Value.Set("bound_value")
				flag.Changed = true
				c.BindFlag("test", flag)
				return c
			}(),
			args:  args{key: "test"},
			want:  &pflag.Flag{Name: "test"},
			want1: true,
		},
		{
			name:  "gets non-existing bound flag",
			c:     NewConfig(),
			args:  args{key: "nonexistent"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.GetBoundFlag(tt.args.key)
			if tt.want1 && got.Name != tt.want.Name {
				t.Errorf("Miniv.GetBoundFlag() got name = %v, want %v", got.Name, tt.want.Name)
			}
			if got1 != tt.want1 {
				t.Errorf("Miniv.GetBoundFlag() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMiniv_GetBoundFlagValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *Miniv
		args  args
		want  any
		want1 bool
	}{
		{
			name: "gets bound flag value when changed",
			c: func() *Miniv {
				c := NewConfig()
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("test", "default", "test flag")
				flag := fs.Lookup("test")
				flag.Value.Set("changed-value")
				flag.Changed = true
				c.BindFlag("test", flag)
				return c
			}(),
			args:  args{key: "test"},
			want:  "changed-value",
			want1: true,
		},
		{
			name: "gets no value when flag not changed",
			c: func() *Miniv {
				c := NewConfig()
				fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
				fs.String("test", "default", "test flag")
				flag := fs.Lookup("test")
				flag.Changed = false
				c.BindFlag("test", flag)
				return c
			}(),
			args:  args{key: "test"},
			want:  nil,
			want1: false,
		},
		{
			name:  "gets no value for non-existing flag",
			c:     NewConfig(),
			args:  args{key: "nonexistent"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.GetBoundFlagValue(tt.args.key)
			if got1 != tt.want1 {
				t.Errorf("Miniv.GetBoundFlagValue() got1 = %v, want %v", got1, tt.want1)
			}
			if got1 {
				if str, ok := got.(interface{ String() string }); ok && str.String() != tt.want {
					t.Errorf("Miniv.GetBoundFlagValue() got = %v, want %v", str.String(), tt.want)
				}
			}
		})
	}
}

func TestMiniv_AutomaticEnv(t *testing.T) {
	tests := []struct {
		name string
		v    *Miniv
	}{
		{
			name: "enables automatic env",
			v:    NewConfig(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.AutomaticEnv()
			if !tt.v.automaticEnvApplied {
				t.Errorf("AutomaticEnv() did not set automaticEnvApplied to true")
			}
		})
	}
}

func TestMiniv_SetEnvVar(t *testing.T) {
	type args struct {
		key    string
		envVar string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets env var",
			c:    NewConfig(),
			args: args{key: "test", envVar: "TEST_VAR"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetEnvVar(tt.args.key, tt.args.envVar)
			expectedKey := strings.ToUpper(tt.args.key)
			if val, exists := tt.c.envVars[expectedKey]; !exists || val != tt.args.envVar {
				t.Errorf("SetEnvVar() did not set env var correctly")
			}
		})
	}
}

func TestMiniv_GetEnvVar(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		c          *Miniv
		args       args
		wantVal    string
		wantExists bool
		setupEnv   func()
		cleanupEnv func()
	}{
		{
			name: "gets env var with automatic env enabled",
			c: func() *Miniv {
				c := NewConfig()
				c.AutomaticEnv()
				return c
			}(),
			args:       args{key: "TEST_VAR"},
			wantVal:    "test_value",
			wantExists: true,
			setupEnv: func() {
				t.Setenv("TEST_VAR", "test_value")
			},
			cleanupEnv: func() {},
		},
		{
			name: "gets no env var when empty and not valid",
			c: func() *Miniv {
				c := NewConfig()
				c.AutomaticEnv()
				return c
			}(),
			args:       args{key: "EMPTY_VAR"},
			wantVal:    "",
			wantExists: false,
			setupEnv: func() {
				t.Setenv("EMPTY_VAR", "")
			},
			cleanupEnv: func() {},
		},
		{
			name: "gets explicitly set env var",
			c: func() *Miniv {
				c := NewConfig()
				c.SetEnvVar("test", "CUSTOM_VAR")
				return c
			}(),
			args:       args{key: "test"},
			wantVal:    "custom_value",
			wantExists: true,
			setupEnv: func() {
				t.Setenv("CUSTOM_VAR", "custom_value")
			},
			cleanupEnv: func() {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv()
			}
			gotVal, gotExists := tt.c.GetEnvVar(tt.args.key)
			if gotVal != tt.wantVal {
				t.Errorf("Miniv.GetEnvVar() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotExists != tt.wantExists {
				t.Errorf("Miniv.GetEnvVar() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
			if tt.cleanupEnv != nil {
				tt.cleanupEnv()
			}
		})
	}
}

func TestMiniv_SetDefault(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		{
			name: "sets default value",
			c:    NewConfig(),
			args: args{key: "test", value: "default_value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.SetDefault(tt.args.key, tt.args.value)
			if val, exists := tt.c.defaults[tt.args.key]; !exists || val != tt.args.value {
				t.Errorf("SetDefault() did not set default correctly")
			}
		})
	}
}

func TestMiniv_GetDefault(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *Miniv
		args  args
		want  any
		want1 bool
	}{
		{
			name: "gets existing default",
			c: func() *Miniv {
				c := NewConfig()
				c.SetDefault("test", "default_value")
				return c
			}(),
			args:  args{key: "test"},
			want:  "default_value",
			want1: true,
		},
		{
			name:  "gets non-existing default",
			c:     NewConfig(),
			args:  args{key: "nonexistent"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.GetDefault(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetDefault() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Miniv.GetDefault() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMiniv_GetConfigValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *Miniv
		args  args
		want  any
		want1 bool
	}{
		{
			name: "gets value from cfgValues",
			c: func() *Miniv {
				c := NewConfig()
				c.cfgValues = map[string]any{"test": "value"}
				return c
			}(),
			args:  args{key: "test"},
			want:  "value",
			want1: true,
		},
		{
			name: "gets value from flatCfgValues",
			c: func() *Miniv {
				c := NewConfig()
				c.flatCfgValues = map[string]any{"flat": "flat_value"}
				return c
			}(),
			args:  args{key: "flat"},
			want:  "flat_value",
			want1: true,
		},
		{
			name:  "gets no value for non-existing key",
			c:     NewConfig(),
			args:  args{key: "nonexistent"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.GetConfigValue(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetConfigValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Miniv.GetConfigValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMiniv_flattenConfigValues(t *testing.T) {
	type args struct {
		prefix  string
		values  map[string]any
		flatMap map[string]any
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.flattenConfigValues(tt.args.prefix, tt.args.values, tt.args.flatMap)
		})
	}
}

func TestMiniv_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		c          *Miniv
		args       args
		want       any
		want1      bool
		setupEnv   func()
		cleanupEnv func()
	}{
		{
			name: "gets set value (highest precedence)",
			c: func() *Miniv {
				c := NewConfig()
				c.SetValue("test", "set_value")
				c.SetDefault("test", "default")
				return c
			}(),
			args:  args{key: "test"},
			want:  "set_value",
			want1: true,
		},
		{
			name: "gets env var value",
			c: func() *Miniv {
				c := NewConfig()
				c.AutomaticEnv()
				return c
			}(),
			args:       args{key: "ENV_TEST"},
			want:       "env_value",
			want1:      true,
			setupEnv:   func() { t.Setenv("ENV_TEST", "env_value") },
			cleanupEnv: func() {},
		},
		{
			name: "gets config value",
			c: func() *Miniv {
				c := NewConfig()
				c.flatCfgValues = map[string]any{"config": "config_value"}
				return c
			}(),
			args:  args{key: "config"},
			want:  "config_value",
			want1: true,
		},
		{
			name: "gets default value (lowest precedence)",
			c: func() *Miniv {
				c := NewConfig()
				c.SetDefault("default", "default_value")
				return c
			}(),
			args:  args{key: "default"},
			want:  "default_value",
			want1: true,
		},
		{
			name:  "gets nothing for non-existing key",
			c:     NewConfig(),
			args:  args{key: "nonexistent"},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupEnv != nil {
				tt.setupEnv()
			}
			got, got1 := tt.c.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Miniv.Get() got1 = %v, want %v", got1, tt.want1)
			}
			if tt.cleanupEnv != nil {
				tt.cleanupEnv()
			}
		})
	}
}

func TestMiniv_GetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want string
	}{
		{
			name: "gets string value",
			c: func() *Miniv {
				c := NewConfig()
				c.SetValue("test", "string_value")
				return c
			}(),
			args: args{key: "test"},
			want: "string_value",
		},
		{
			name: "gets empty string for non-string value",
			c: func() *Miniv {
				c := NewConfig()
				c.SetValue("test", 123)
				return c
			}(),
			args: args{key: "test"},
			want: "123",
		},
		{
			name: "gets empty string for non-existing key",
			c:    NewConfig(),
			args: args{key: "nonexistent"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetString(tt.args.key); got != tt.want {
				t.Errorf("Miniv.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetStringSlice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetStringSlice(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetInt64(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetInt64(tt.args.key); got != tt.want {
				t.Errorf("Miniv.GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetInt64Slice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetInt64Slice(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetInt64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetFloat64(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetFloat64(tt.args.key); got != tt.want {
				t.Errorf("Miniv.GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetFloat64Slice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetFloat64Slice(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetFloat64Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetBool(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetBool(tt.args.key); got != tt.want {
				t.Errorf("Miniv.GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_GetBoolSlice(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		c    *Miniv
		args args
		want []bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetBoolSlice(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Miniv.GetBoolSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_ConfigFileUsed(t *testing.T) {
	tests := []struct {
		name string
		c    *Miniv
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ConfigFileUsed(); got != tt.want {
				t.Errorf("Miniv.ConfigFileUsed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiniv_ReadInConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *Miniv
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ReadInConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Miniv.ReadInConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMiniv_WriteConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *Miniv
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.WriteConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Miniv.WriteConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMiniv_WriteConfigAs(t *testing.T) {
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name    string
		c       *Miniv
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.WriteConfigAs(tt.args.cfgFile); (err != nil) != tt.wantErr {
				t.Errorf("Miniv.WriteConfigAs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMiniv_SafeWriteConfig(t *testing.T) {
	tests := []struct {
		name    string
		c       *Miniv
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SafeWriteConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Miniv.SafeWriteConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMiniv_SafeWriteConfigAs(t *testing.T) {
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name    string
		c       *Miniv
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SafeWriteConfigAs(tt.args.cfgFile); (err != nil) != tt.wantErr {
				t.Errorf("Miniv.SafeWriteConfigAs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
