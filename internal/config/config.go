package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"go-blog-admin/internal/config/consts"
	"go-blog-admin/internal/util/utilconfig"
	xlog "go-blog-admin/internal/util/utillog"
	"os"
	"slices"
	"strconv"
	"strings"
)

var (
	AppVersion  = ""
	AppCommit   = ""
	AppDate     = ""
	ShortCommit = ""
)

func dumpVersionAndExitIf() {

	if CmdLine.Version {
		fmt.Printf("Version: %s\n", AppVersion)
		fmt.Printf("Commit: %s\n", AppCommit)
		fmt.Printf("Date: %s\n", AppDate)
		//
		os.Exit(0)
	}

}

type CmdLineConfig struct {
	Config  string
	CertDir string

	Env     string
	Name    string
	Version bool

	SysAPIKey string
	Listen    string
	ListenTLS string
	ListenSys string

	DumpConfig bool
}

const (
	envDevelopment = "development"
	envTesting     = "testing"
	envStaging     = "staging"
	envProduction  = "production"
)

var envNames = []string{
	envDevelopment, envTesting, envStaging, envProduction,
}

var CmdLine = CmdLineConfig{}

// ReadFlags read app flags
func ReadFlags() {
	_ = os.Args
	flag.StringVar(&CmdLine.Config, "config", "", "Path to dir with config files")
	flag.StringVar(&CmdLine.CertDir, "cert-dir", "", "Path to dir with cert files")
	flag.StringVar(&CmdLine.SysAPIKey, "sys-api-key", "", "Sys api key")
	flag.StringVar(&CmdLine.Listen, "listen", "", "Listen")
	flag.StringVar(&CmdLine.ListenTLS, "listen-tls", "", "Listen TLS")
	flag.StringVar(&CmdLine.ListenSys, "listen-sys", "", "Listen sys")

	flag.StringVar(&CmdLine.Env, "env", "", "Environment: development, testing, staging, production")
	flag.StringVar(&CmdLine.Name, "name", "", "App name")

	flag.BoolVar(&CmdLine.Version, "version", false, "App version")

	flag.BoolVar(&CmdLine.DumpConfig, "dump-config", false, "Dump Config")

	flag.Parse() // dont use from init()

	dumpVersionAndExitIf()
}

type envReader struct {
	envError error
	prefix   string
}

func NewEnvReader() envReader {
	return envReader{prefix: "app_"}
}
func (x *envReader) String(p *string, name string, cmdValue *string) {
	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive
	if cmdValue != nil && *cmdValue != "" {
		xlog.Info("Reading %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}

	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			xlog.Info("Reading %q value from env: %v = %v", name, envName, envValue)
			*p = envValue
			return
		}
	}

}

func (x *envReader) Bool(p *bool, name string, cmdValue *bool) {

	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive

	if cmdValue != nil && *cmdValue {
		xlog.Info("Reading %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}
	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			xlog.Info("Reading %q value from env: %v = %v", name, envName, envValue)
			*p = envValue == "1" || envValue == "true"
			return
		}
	}
}

func (x *envReader) Int(p *int, name string, cmdValue *int) {

	envName := strings.ToUpper(x.prefix + name) // *nix case-sensitive

	if cmdValue != nil && *cmdValue != 0 {
		xlog.Info("Reading %q value from cmd: %v", name, *cmdValue)
		*p = *cmdValue
		return
	}
	if envName != "" {
		envValue := os.Getenv(envName)
		if envValue != "" {
			xlog.Info("Reading %q value from env: %v = %v", name, envName, envValue)

			if v, err := strconv.Atoi(envValue); err == nil {
				*p = v
			} else {
				x.envError = err
			}

		}
	}

}

type Database struct {
	Dialect   string `json:"dialect"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	Name      string `json:"name"`
	User      string `json:"user"`
	Password  string `json:"password"`
	MaxOpen   int    `json:"max_open"`
	MaxIdle   int    `json:"max_idle"`
	IdleTime  int    `json:"idle_time"`
	Migration bool   `json:"migration"`
	Debug     bool   `json:"debug"`
}

// type AppConfigLogger struct {
// 	Level int `json:"level"`
// }

type AppConfigVaultKey struct {
	ID      string `json:"id"`
	AuthKey string `json:"auth_key"` // for user auth (used by cluster apps)
}

func (x AppConfigVaultKey) IsEmpty() bool {
	return x.ID == "" || x.AuthKey == ""
}

type AppConfigVault struct {
	Keys []AppConfigVaultKey `json:"keys"` // keychain
}

type AppConfigIdentity struct {
	AuthTokenIssuer string `json:"auth_token_issuer"` // default "auth"

}

type AppConfigBlog struct {
}

func (x *AppConfigBlog) Validate() error {

	return nil
}

type AppConfigLang struct {
	Langs []string `json:"langs"`
}

type AppConfigAssets struct {
	GlobalVersion   string
	AssetsPublicURL string
}
type AppConfigMod struct {
	Name  string `json:"-"`
	Env   string `json:"env"` // prod||'' dev stage
	Debug bool   `json:"-"`
	Title string `json:"title"`

	ConfigPath []string `json:"-"` // []string{".", os.Getenv("APP_CONFIG"), flagAppConfig}
}

type AppConfigHTTPTransport struct {
	MaxIdleConns        int `json:"max_idle_conns,omitempty"`
	MaxIdleConnsPerHost int `json:"max_idle_conns_per_host,omitempty"`
	IdleConnTimeout     int `json:"idle_conn_timeout,omitempty"`
	MaxConnsPerHost     int `json:"max_conns_per_host,omitempty"`
}
type AppConfigHTTPServer struct {
	AccessLog bool `json:"access_log"`

	RateLimit     float64 `json:"rate_limit"`
	RateBurst     int     `json:"rate_burst"`
	Listen        string  `json:"listen"`
	ListenTLS     string  `json:"listen_tls"`
	AutoTLS       bool    `json:"auto_tls"`
	RedirectHTTPS bool    `json:"redirect_https"`
	RedirectWWW   bool    `json:"redirect_www"`

	CertDir string `json:"cert_dir"`

	ReadTimeout       int `json:"read_timeout,omitempty"`        // 5 to 30 seconds
	WriteTimeout      int `json:"write_timeout,omitempty"`       // 10 to 30 seconds, WriteTimeout > ReadTimeout
	IdleTimeout       int `json:"idle_timeout,omitempty"`        // 60 to 120 seconds
	ReadHeaderTimeout int `json:"read_header_timeout,omitempty"` // default get from ReadTimeout

	SysMetrics bool   `json:"sys_metrics"` //
	SysAPIKey  string `json:"sys_api_key"`
	ListenSys  string `json:"listen_sys"`
}
type AppConfig struct {
	AppConfigMod

	// Logger AppConfigLogger `json:"logger"`

	Vault AppConfigVault `json:"vault"`

	Identity AppConfigIdentity `json:"identity"`

	Blog AppConfigBlog `json:"blog"`

	DB    Database `json:"database"`
	Redis Database `json:"redis"`

	Lang AppConfigLang `json:"lang"`

	Assets AppConfigAssets `json:"assets"`

	HTTPTransport AppConfigHTTPTransport `json:"http_transport"`

	HTTPServer AppConfigHTTPServer `json:"http_server"`
}

func NewAppConfig() *AppConfig {

	res := &AppConfig{

		Lang: AppConfigLang{Langs: []string{"en"}},

		// Logger: AppConfigLogger{
		// 	Level: consts.LogLevelWarn,
		// },

		Vault: AppConfigVault{
			Keys: []AppConfigVaultKey{},
		},

		Identity: AppConfigIdentity{
			AuthTokenIssuer: "auth",
		},

		Blog: AppConfigBlog{},

		DB: Database{
			Dialect:  "postgres",
			Host:     "127.0.0.1",
			Port:     "5432",
			Name:     "postgres",
			User:     "postgres",
			Password: "postgres",
			MaxOpen:  0,
			MaxIdle:  0,
			IdleTime: 0,
		},
		Redis: Database{
			Host:     "127.0.0.1",
			Port:     "6379",
			Name:     "redis",
			User:     "redis",
			Password: "redis",
		},

		AppConfigMod: AppConfigMod{
			Name:       consts.AppName,
			ConfigPath: []string{"."},
			Title:      "",
			Env:        "production",
			Debug:      false,
		},

		Assets: AppConfigAssets{
			GlobalVersion:   "v1",
			AssetsPublicURL: "",
		},

		HTTPTransport: AppConfigHTTPTransport{},

		HTTPServer: AppConfigHTTPServer{
			ReadTimeout:  0,
			WriteTimeout: 0,
			IdleTimeout:  0,

			RateLimit: 0,
			RateBurst: 0,

			Listen: "127.0.0.1:18180",
			// ListenTLS: "127.0.0.1:10283",

			CertDir: "",

			SysAPIKey: "",
		},
	}

	return res
}

func (x *AppConfig) readEnvName() error {
	reader := NewEnvReader()
	// APP_ENV -env
	reader.String(&x.Env, "env", &CmdLine.Env)
	reader.String(&x.Name, "name", &CmdLine.Name)

	if err := x.validateEnv(); err != nil {
		return err
	}

	configPath := slices.Concat(strings.Split(os.Getenv("APP_CONFIG"), ";"), strings.Split(CmdLine.Config, ";"))
	configPath = slices.Compact(configPath)
	configPath = slices.DeleteFunc(
		configPath,
		func(x string) bool {
			return x == ""
		},
	)

	for i := 0; i < len(configPath); i++ {
		configPath[i] += "/" + x.Name
	}

	// if len(configPath) == 0 {
	// 	configPath = []string{"."} // default
	// }

	if len(configPath) == 0 {
		xlog.Warn("Config path is empty")
	} else {
		xlog.Info("Config path: %v", configPath)
	}

	x.ConfigPath = configPath

	return nil
}
func (x *AppConfig) readEnvVar() error {
	reader := NewEnvReader()

	// Assets configuration
	reader.String(&x.Assets.GlobalVersion, "global_version", nil)

	// Database configuration

	reader.String(&x.DB.Dialect, "db_dialect", nil)
	reader.String(&x.DB.Host, "db_host", nil)
	reader.String(&x.DB.Port, "db_port", nil)
	reader.String(&x.DB.Name, "db_name", nil)
	reader.String(&x.DB.User, "db_user", nil)
	reader.String(&x.DB.Password, "db_password", nil)
	reader.Int(&x.DB.MaxOpen, "db_max_open", nil)
	reader.Int(&x.DB.MaxIdle, "db_max_idle", nil)
	reader.Int(&x.DB.IdleTime, "db_idle_time", nil)
	reader.Bool(&x.DB.Migration, "db_migration", nil)

	// General configuration
	reader.String(&x.Title, "title", nil)

	reader.String(&x.HTTPServer.CertDir, "cert_dir", &CmdLine.CertDir)

	reader.String(&x.HTTPServer.Listen, "listen", &CmdLine.Listen)
	reader.String(&x.HTTPServer.ListenTLS, "listen_tls", &CmdLine.ListenTLS)
	reader.String(&x.HTTPServer.ListenSys, "listen_sys", &CmdLine.ListenSys)

	reader.String(&x.HTTPServer.SysAPIKey, "sys_api_key", &CmdLine.SysAPIKey)

	if reader.envError != nil {
		return reader.envError
	}

	return nil
}

func (x *AppConfig) validateEnv() error {

	if x.Env == "" {
		x.Env = envProduction
	}

	x.Debug = x.Env == envDevelopment

	if !slices.Contains(envNames, x.Env) {
		xlog.Warn("Non-standart env name: %v", x.Env)
	}

	return nil

}
func (x AppConfig) validate() error {

	if x.HTTPServer.Listen == "" && x.HTTPServer.ListenTLS == "" {
		return fmt.Errorf("socket Listen and ListenTLS are empty")
	}

	if err := x.Blog.Validate(); err != nil {
		return fmt.Errorf("error on .Blog.Validate: %v", err)
	}
	return nil
}

type AppConfigSource struct {
	config *AppConfig
}

func MustNewAppConfigSource() *AppConfigSource {

	res := &AppConfigSource{}

	err := res.Load() // init

	if err != nil {
		panic(err) // must
	}

	return res

}

func (x *AppConfigSource) Load() error {

	res := NewAppConfig()

	{
		err := res.readEnvName()
		if err != nil {
			return err
		}
	}

	{

		for i := 0; i < len(res.ConfigPath); i++ {

			dir := res.ConfigPath[i]
			fileName := fmt.Sprintf("config.%s.json", res.Env)

			xlog.Info("Loading config from: %v", dir)

			err := utilconfig.LoadConfig(res /*pointer*/, dir, fileName)

			if err != nil {
				return err
			}

		}

	}

	{
		err := res.readEnvVar()
		if err != nil {
			return err
		}

	}

	{
		err := res.validate()
		if err != nil {
			return err
		}
	}

	xlog.Info("Config loaded: Name=%v Env=%v Debug=%v ", res.Name, res.Env, res.Debug)

	x.config = res

	if CmdLine.DumpConfig {
		data, _ := json.MarshalIndent(res, "", " ")
		fmt.Println(string(data))
	}

	return nil
}

func (x *AppConfigSource) Config() *AppConfig {

	return x.config

}
