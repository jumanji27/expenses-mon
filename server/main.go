package main

import (
    "fmt"
    // "reflect"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "money_mon/server/model"
    "money_mon/server/router"
)


type Main struct {}

func main() {
    martini_app := martini.Classic()
    martini_app.Use(render.Renderer())

    app := Main{}
    model.init()
    router.init(martini_app)

    fmt.Printf("App starting!\n")

    martini_app.Run()
}
