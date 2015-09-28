package main

import (
  // "fmt"
  // "reflect"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"

  "expense_mon/server/router"
  "expense_mon/server/helpers"
)


func main() {
  martini_app := martini.Classic()

  martini_app.Use(render.Renderer())
  martini_app.Use(
    render.Renderer(
      render.Options{
        Directory: "server/tmpl",
      },
    ),
  )
  martini_app.Use(
    martini.Static("client/public"),
  )

  router := router.Main{}
  router.Init(martini_app)

  helpers := helpers.Main{}
  helpers.LogSimpleMessage("App starting!")

  martini_app.Run()
}
