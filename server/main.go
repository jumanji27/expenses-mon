package main

import (
    "fmt"
    // "reflect"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "money_mon/server/router"
)


func main() {
    martini_app := martini.Classic()
    martini_app.Use(render.Renderer())

    router := router.Main{}
    router.Init(martini_app)

    fmt.Printf("App starting!\n")

    martini_app.Run()
}
