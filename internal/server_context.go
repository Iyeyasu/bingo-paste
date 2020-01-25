package server

import (
	"database/sql"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/view"
	"github.com/julienschmidt/httprouter"
)

type ServerContext struct {
	Config   *config.Config
	Database *sql.DB
	Renderer *view.TemplateRenderer
	Router   *httprouter.Router
}
