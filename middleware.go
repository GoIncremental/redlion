package main

import (
	"net/http"
	"os"

	"github.com/goincremental/dal"
	"github.com/goincremental/redlion/models"
	"github.com/goincremental/web"
)

func Mongo(connection dal.Connection) web.Middleware {
	database := os.Getenv("MONGO_NAME")

	return web.MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		conn := connection.Clone()
		defer conn.Close()
		db := conn.DB(database)
		web.SetDb(r, db)
		next(rw, r)
	})
}

func Site() web.Middleware {
	site := &models.Site{
		Env: os.Getenv("APP_ENV"),
	}
	return web.MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		SetSite(r, site)
		next(rw, r)
	})
}

//Render creates a default web.Renderer and then returns a middleware function
//that ensures this is placed on the request context for each request
func Render() web.Middleware {
	renderer := web.NewRenderer()

	return web.MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		web.SetRenderer(r, renderer)
		next(rw, r)
	})
}
