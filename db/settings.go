package db

import (
	"database/sql"
	"fmt"
	"log"
)

var businessName string = "BusinessName"
var businessPhone string = "BusinessPhone"
var businessAddr string = "BusinessAddr"
var businessCity string = "BusinessCity"
var businessState string = "BusinessState"
var businessZip string = "BusinessZip"
var invoiceColor string = "InvoiceColor"
var invoiceFooter string = "InvoiceFooter"
var emailBccs string = "EmailBccs"

type IndigoSettings struct {
	BusinessName  string
	BusinessPhone string

	BusinessAddr  string
	BusinessCity  string
	BusinessState string
	BusinessZip   string

	InvoiceColor  string
	InvoiceFooter string

	EmailBccs string
}

func (s IndigoSettings) ToDict() map[string]string {
	res := make(map[string]string, 0)

	res[businessName] = s.BusinessName
	res[businessPhone] = s.BusinessPhone
	res[businessAddr] = s.BusinessAddr
	res[businessCity] = s.BusinessCity
	res[businessState] = s.BusinessState
	res[businessZip] = s.BusinessZip
	res[invoiceColor] = s.InvoiceColor
	res[invoiceFooter] = s.InvoiceFooter
	res[emailBccs] = s.EmailBccs

	return res
}

func SettingsFromDict(dict map[string]string) IndigoSettings {
	return IndigoSettings{
		BusinessName:  dict[businessName],
		BusinessPhone: dict[businessPhone],
		BusinessAddr:  dict[businessAddr],
		BusinessCity:  dict[businessCity],
		BusinessState: dict[businessState],
		BusinessZip:   dict[businessZip],
		InvoiceColor:  dict[invoiceColor],
		InvoiceFooter: dict[invoiceFooter],
		EmailBccs:     dict[emailBccs],
	}
}

type SettingsTable struct {
	db *sql.DB
}

var settingsInstance *SettingsTable

func InitSettingsTable(db *sql.DB) *SettingsTable {
	log.Println("Initializing Settings table.")
	if settingsInstance != nil {
		log.Println("Settings table already initialized, returning existing instance.")
		return settingsInstance
	}

	log.Println("Ensuring Settings table exists.")
	settingsInstance = &SettingsTable{db: db}
	_, err := settingsInstance.db.Exec(`
        CREATE TABLE IF NOT EXISTS settings(
            key TEXT PRIMARY KEY,
			value TEXT
        );
    `)

	if err != nil {
		log.Fatal("Error while creating Settings table.", err)
	}

	log.Println("Settings table successfully initialized.")
	return settingsInstance
}

func (h *SettingsTable) Get(key string) (string, error) {
	row := h.db.QueryRow(`
		SELECT value,
		FROM settings 
		WHERE key = (?);`, key)

	var value string
	err := row.Scan(&value)

	return value, err
}

func (h *SettingsTable) List() (settings map[string]string, error error) {
	settings = make(map[string]string, 0)

	rows, error := h.db.Query(`
		SELECT key, value
		FROM settings;`)

	if error != nil {
		return
	}

	for rows.Next() {
		error = rows.Err()
		if error != nil {
			return
		}

		var key, value string
		if error = rows.Scan(&key, &value); error != nil {
			return
		}

		settings[key] = value
	}

	return settings, nil
}

func (h *SettingsTable) SetMany(settings map[string]string) error {
	log.Println("Updating settings")

	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	stmnt, err := tx.Prepare(`
		REPLACE INTO settings(key, value)
		VALUES((?), (?));`)
	fmt.Println("1a")

	if err != nil {
		return err
	}

	defer stmnt.Close()

	for key, value := range settings {
		_, err = stmnt.Exec(key, value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Println("Successfully updated settings")

	return nil
}

func (h *SettingsTable) GetAll() (sett IndigoSettings, error error) {
	dict, error := h.List()
	if error != nil {
		return
	}

	sett = SettingsFromDict(dict)
	return
}

func (h *SettingsTable) SetAll(sett IndigoSettings) error {
	dict := sett.ToDict()
	err := h.SetMany(dict)
	return err
}
