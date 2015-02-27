package main


import (
    "fmt"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    // "labix.org/v2/mgo"
)


type PagesData struct {
    Number      int
    Description string
    Status      string
}


func setData() PagesData {
    data := PagesData{
        Number:      123,
        Description: "Page number",
        Status:      "Done!",
    }
    return data
}


func main() {
    app := martini.Classic()
    app.Use(render.Renderer())

    app.Get("/api/v1/index", func(render render.Render, params martini.Params) {
        render.JSON(200, setData())
    })

    fmt.Printf("App starting!\n")

    app.Run()
}