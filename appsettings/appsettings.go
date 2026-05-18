package appsettings

import (
	"os"

	"github.com/BurntSushi/toml"
)

const SettingsPath = "indigo.toml"

type SmtpSettings struct {
	Host     string
	Port     int
	From     string
	PassFile string
}

type AppSettings struct {
	DbPath string
	Smtp   SmtpSettings
}

func GetSettings() (settings AppSettings, error error) {
	cfg := os.Getenv("INDIGO_CONFIG")
	if cfg == "" {
		cfg = SettingsPath
	}

	bytes, error := os.ReadFile(cfg)
	if error != nil {
		return
	}

	toml.Decode(string(bytes), &settings)
	return
}
