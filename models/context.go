package models

import (
	"github.com/goincremental/dal"
	"github.com/goincremental/web"
)

type Context struct {
	DB     dal.Database
	Render web.Renderer
	Site   *Site
	SiteID *dal.ObjectID
}
