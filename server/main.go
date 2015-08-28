package main

import (
    // "reflect"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "money_mon/server/router"
    "money_mon/server/helpers"
)


func main() {
    martini_app := martini.Classic()
    martini_app.Use(render.Renderer())

    router := router.Main{}
    router.Init(martini_app)

    helpers := helpers.Main{}
    helpers.LogSimpleMessage("App starting!")

    martini_app.Run()
}
