package api

import (
	"gihub.com/jastribl/balancedot/config"
	"gihub.com/jastribl/balancedot/repos"
	"gihub.com/jastribl/balancedot/services/auth"
	"gorm.io/gorm"
)

// App contains API handler methods
type App struct {
	db     *gorm.DB
	config *config.Config
	Auth   *auth.Auth
}

// NewApp returns a new App
func NewApp(db *gorm.DB, config *config.Config) (*App, error) {
	newAuth := auth.NewAPIAuth(db, "todo-cookie-name")
	return &App{
		db:     db,
		config: config,
		Auth:   newauth,
	}, nil
}

func (m *App) genericGetAll(
	w ResponseWriter,
	r *Request,
	repo *gorm.DB,
	typeRef interface{},
	options *repos.GetAllOfOptions,
) WriterResponse {
	items, err := repos.NewGenericRepo(repo).GetAllOf(typeRef, options)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(items)
}

func (m *App) genricRawFindAll(
	w ResponseWriter,
	r *Request,
	repo *gorm.DB,
	typeRef interface{},
	query string,
	params ...interface{},
) WriterResponse {
	items, err := repos.NewGenericRepo(repo).GetAllOfRaw(typeRef, query, params...)
	if err != nil {
		return w.SendUnexpectedError(err)
	}

	return w.SendResponse(items)
}

func (m *App) genericGetByUUID(
	w ResponseWriter,
	r *Request,
	repo *gorm.DB,
	out interface{},
	uuid string,
) WriterResponse {
	err := repos.NewGenericRepo(repo).GetByUUID(out, uuid)
	if err != nil {
		return w.SendUnexpectedError(err)
	}
	return w.SendResponse(out)
}
