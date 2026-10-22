package routes

import (
	"github.com/aosanya/mwanachama-backend-shared/dispatch"
	"github.com/aosanya/mwanachama-backend-shared/httpwire"

	engage "github.com/aosanya/mwanachama-backend-engage"
)

type Route = httpwire.Route

type Mount = dispatch.Mount

var Sentinels = map[string]error{
	"ErrPostNotFound": engage.ErrPostNotFound,
	"ErrInvalidPost":  engage.ErrInvalidPost,
	"ErrNoPoster":     engage.ErrNoPoster,
}

var AnonymousActions = []string{}

var Table = dispatch.NewTable(engage.Operations(), Sentinels, AnonymousActions...)

func Build(pm engage.PostManager) ([]Route, error) { return Table.Build(pm, Mount{}) }

func BuildFor(pm engage.PostManager, m Mount) ([]Route, error) { return Table.Build(pm, m) }

func RoutesFor(pm engage.PostManager, m Mount) []Route { return Table.Routes(pm, m) }

func Split(pm engage.PostManager, m Mount) dispatch.Split { return Table.Split(pm, m) }
