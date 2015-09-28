package router

import (
    // "fmt"
    // "reflect"
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "expense_mon/server/model"
)


type Main struct {}


func (self *Main) Init(app *martini.ClassicMartini) {
    const (
        http_success = 200
    )

    app.Get(
        "/**",
        func(render render.Render) {
            render.HTML(http_success, "index", nil)
        },
    )

    const api_base_url = "/api/v1"

    app.Post(
        api_base_url,
        func(render render.Render) {
            render.JSON(
                http_success,
                map[string]interface{}{
                    "success": map[string]interface{}{"greeting": "Hello, I'm your API!"},
                    "error": nil,
                },
            )
        },
    )

    model := model.Main{}
    model.Init()

    api_url := api_base_url
    handler := "/get"
    api_url += handler

    app.Post(
        api_url,
        func(render render.Render) {
            render.JSON(http_success, model.Get())
        },
    )

    api_url = api_base_url
    handler = "/set"
    api_url += handler

    app.Post(
        api_url,
        func(res *http.Request, render render.Render) {
            render.JSON(http_success, model.Set(res))
        },
    )
}