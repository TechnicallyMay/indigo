package appsettings

import (
	"os"

	"github.com/BurntSushi/toml"
)

var SettingsPath = "indigo.toml"

type SmtpSettings struct {
	Host     string
	Port     int
	From     string
	PassFile string
}

type AppSettings struct {
	Smtp SmtpSettings
}

func GetSettings() (settings AppSettings, error error) {
	bytes, error := os.ReadFile(SettingsPath)
	if error != nil {
		return
	}

	toml.Decode(string(bytes), &settings)
	return
}
