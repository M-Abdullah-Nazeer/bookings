package config

import (
	"log"
	"text/template"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
