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
    "/",
    func(render render.Render) {
      render.HTML(httpSuccess, "main", nil)
    },
  )

  const APIBaseURL = "/api/v1"

  app.Post(
    APIBaseURL,
    func(render render.Render) {
      render.JSON(
        httpSuccess,
        map[string]interface{}{
          "success": map[string]interface{}{"greeting": "Hello, I'm your API!"},
        },
      )
    },
  )

  expensesModel := expensesModel.Main{}
  expensesModel.Init()

  APIURL := APIBaseURL
  handler := "/get"
  APIURL += handler

  app.Post(
    APIURL,
    func(render render.Render) {
      render.JSON(httpSuccess, expensesModel.GetHandler())
    },
  )

  APIURL = APIBaseURL
  handler = "/set"
  APIURL += handler

  app.Post(
    APIURL,
    func(res *http.Request, render render.Render) {
      render.JSON(
        httpSuccess,
        expensesModel.SetHandler(res),
      )
    },
  )
}