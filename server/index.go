package main


import (
    "fmt"
    "log"
    // "reflect"

    "encoding/json"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
)

type Index struct {
    MongoCollection *mgo.Collection
    Data
}

type Data struct {
    Value [1][3][12][5]int
}


func (self Index) db_init() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }

    defer session.Close()

    self.MongoCollection = session.DB("test").C("test")
}

func (self Index) db_get() string {
    raw_result := Data{}
    self.MongoCollection.Find(nil).One(&raw_result)

    result, err := json.Marshal(raw_result.Value)
    if err != nil {
        log.Fatal(err)
    }

    return string(result)
}

func(self Index) db_set() string {
    // var local_data [1][3][12][5]int

    // for i := 0; i < 3; i++ {
    //     for j := 0; j < 12; j++ {
    //         for k := 0; k < 5; k++ {
    //             local_data[0][i][j][k] = k + 1
    //         }
    //     }
    // }

    // data := self.Data{Value: local_data}

    // err = collection.Insert(data)
    // if err != nil {
    //     log.Fatal(err)
    // }
    return "test"
}


func (self Index) route(martini_app *martini.ClassicMartini) {
    martini_app.Get("/", func(render render.Render) {
        render.JSON(200, map[string]interface{}{
            "success": map[string]interface{}{"greeting": "Hello, I'm your API!"},
            "error": nil,
        })
    })

    martini_app.Post("/api/v1/get", func(render render.Render) {
        render.JSON(200, self.db_get())
    })

    martini_app.Post("/api/v1/set", func(render render.Render) {
        render.JSON(200, self.db_set())
    })
}


func main() {
    martini_app := martini.Classic()
    martini_app.Use(render.Renderer())

    app := Index{}
    app.route(martini_app)

    db_init()

    fmt.Printf("App starting!\n")

    martini_app.Run()
}
