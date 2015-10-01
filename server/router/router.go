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
    httpSuccess = 200
  )

  app.Get(
    "/**",
    func(render render.Render) {
      render.HTML(httpSuccess, "index", nil)
    },
  )

  const apiBaseURL = "/api/v1"

  app.Post(
    apiBaseURL,
    func(render render.Render) {
      render.JSON(
        httpSuccess,
        map[string]interface{}{
          "success": map[string]interface{}{"greeting": "Hello, I'm your API!"},
          "error": nil,
        },
      )
    },
  )

  model := model.Main{}
  model.Init()

  apiURL := apiBaseURL
  handler := "/get"
  apiURL += handler

  app.Post(
    apiURL,
    func(render render.Render) {
      render.JSON(httpSuccess, model.Get())
    },
  )

  apiURL = apiBaseURL
  handler = "/set"
  apiURL += handler

  app.Post(
    apiURL,
    func(res *http.Request, render render.Render) {
      render.JSON(httpSuccess, model.Set(res))
    },
  )

  apiURL = apiBaseURL
  handler = "/remove"
  apiURL += handler

  app.Post(
    apiURL,
    func(res *http.Request, render render.Render) {
      render.JSON(httpSuccess, model.Remove(res))
    },
  )
}