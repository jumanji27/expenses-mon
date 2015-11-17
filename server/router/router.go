package router

import (
  // "fmt"
  // "reflect"
  "net/http"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"

  "expenses-mon/server/models/expenses"
)


type Main struct {}


func (self *Main) Init(app *martini.ClassicMartini) {
  const (
    httpSuccess = 200
  )

  app.Get(
    "/**",
    func(render render.Render) {
      render.HTML(httpSuccess, "main", nil)
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

  expensesModel := expensesModel.Main{}
  expensesModel.Init()

  apiURL := apiBaseURL
  handler := "/get"
  apiURL += handler

  app.Post(
    apiURL,
    func(render render.Render) {
      render.JSON(httpSuccess, expensesModel.GetHandler())
    },
  )

  apiURL = apiBaseURL
  handler = "/set"
  apiURL += handler

  app.Post(
    apiURL,
    func(res *http.Request, render render.Render) {
      render.JSON(
        httpSuccess,
        expensesModel.SetHandler(res),
      )
    },
  )
}