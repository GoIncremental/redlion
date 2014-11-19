package main

import (
	"log"
	"net/http"
	"os"

	"github.com/goincremental/web"
)

type contextKey int

const site contextKey = 3

//Render creates a default web.Renderer and then returns a middleware function
//that ensures this is placed on the request context for each request
func Render() web.Middleware {
	renderer := web.NewRenderer()

	return web.MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		web.SetRenderer(r, renderer)
		next(rw, r)
	})
}

type Site struct {
	Env string
}

func SetSite(r *http.Request, val *Site) {
	web.SetContext(r, site, val)
}

func GetSite(r *http.Request) *Site {
	if rv := web.GetContext(r, site); rv != nil {
		return rv.(*Site)
	}
	return nil
}

func SiteMW() web.Middleware {
	site := &Site{
		Env: os.Getenv("APP_ENV"),
	}
	return web.MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		SetSite(r, site)
		next(rw, r)
	})
}

type context struct {
	render web.Renderer
	Site   *Site
}

func GetRequestContext(req *http.Request) (ctx *context, err error) {
	r := web.GetRenderer(req)
	s := GetSite(req)
	ctx = &context{
		render: r,
		Site:   s,
	}
	return
}

func getContextAndRender(template string, w http.ResponseWriter, req *http.Request) {
	c, err := GetRequestContext(req)
	if err != nil {
		log.Printf("Error getting request context: %s\n", err)
	}
	c.render.HTML(w, http.StatusOK, template, c)
}

//Home renders the default html home page
func Home(w http.ResponseWriter, req *http.Request) {
	getContextAndRender("home", w, req)
}

func main() {
	web.LoadEnv()

	port := "3200"
	p := os.Getenv("PORT")
	if p != "" {
		port = ":" + p
	}

	router := web.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")

	server := web.NewServer()
	server.Use(web.Gzip())
	server.Use(Render())
	server.Use(SiteMW())
	server.UseHandler(router)
	server.Run(port)
}
