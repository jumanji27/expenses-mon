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
  martiniApp := martini.Classic()

  martiniApp.Use(render.Renderer())
  martiniApp.Use(
    render.Renderer(
      render.Options{
        Directory: "server/tmpl",
      },
    ),
  )
  martiniApp.Use(
    martini.Static("client/public"),
  )

  router := router.Main{}
  router.Init(martiniApp)

  helpers := helpers.Main{}
  helpers.LogSimpleMessage("App starting!")

  martiniApp.Run()
}
