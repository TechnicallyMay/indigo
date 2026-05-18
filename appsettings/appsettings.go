package appsettings

import (
	"os"

	"github.com/BurntSushi/toml"
)

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

func GetSettings(settingsPath string) (settings AppSettings, error error) {
	bytes, error := os.ReadFile(settingsPath)
	if error != nil {
		return
	}

	toml.Decode(string(bytes), &settings)
	return
}
