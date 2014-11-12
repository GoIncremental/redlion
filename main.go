package main

import (
	"log"
	"net/http"
	"os"

	"github.com/goincremental/web"
)

//Render creates a default web.Renderer and then returns a middleware function
//that ensures this is placed on the request context for each request
func Render() web.Middleware {
	renderer := web.NewRenderer()

	return web.MiddlewareFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		web.SetRenderer(r, renderer)
		next(rw, r)
	})
}

type context struct {
	render web.Renderer
}

func GetRequestContext(req *http.Request) (ctx *context, err error) {
	r := web.GetRenderer(req)

	ctx = &context{
		render: r,
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
	server.UseHandler(router)
	server.Run(port)
}
