package main

import (
	"net/http"

	"github.com/goincremental/redlion/models"
	"github.com/goincremental/web"
)

type contextKey int

const site contextKey = 3

func SetSite(r *http.Request, val *models.Site) {
	web.SetContext(r, site, val)
}

func GetSite(r *http.Request) *models.Site {
	if rv := web.GetContext(r, site); rv != nil {
		return rv.(*models.Site)
	}
	return nil
}

func GetRequestContext(req *http.Request) (ctx *models.Context, err error) {
	db := web.GetDb(req)
	r := web.GetRenderer(req)
	siteID, err := models.GetSiteID(db)
	s := GetSite(req)
	ctx = &models.Context{
		DB:     db,
		Render: r,
		Site:   s,
		SiteID: siteID,
	}
	return
}
