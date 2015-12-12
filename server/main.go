package main

import (
  // "fmt"
  // "reflect"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"

  "expenses-mon/server/router"
  "expenses-mon/server/helpers"
)


func main() {
  martiniApp := martini.Classic()

  martiniApp.Use(render.Renderer())
  martiniApp.Use(
    render.Renderer(
      render.Options{
        Directory: "server/views",
      },
    ),
  )
  martiniApp.Use(
    martini.Static("client/public"),
  )

  router := router.Main{}
  router.Init(martiniApp)

  helpers := helpers.Main{}
  helpers.CreateEvent("Log", "EM starting!")

  martiniApp.RunOnAddr(":3000")
  martiniApp.Run()
}
