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


type Data struct {
    Value [1][3][12][5]int
}

func db_init() string {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }

    defer session.Close()

    collection := session.DB("test").C("test")

    var local_data [1][3][12][5]int

    for i := 0; i < 3; i++ {
        for j := 0; j < 12; j++ {
            for k := 0; k < 5; k++ {
                local_data[0][i][j][k] = k + 1
            }
        }
    }

    data := Data{Value: local_data}

    err = collection.Insert(data)
    if err != nil {
        log.Fatal(err)
    }

    raw_result := Data{}
    err = collection.Find(nil).One(&raw_result)
    if err != nil {
        log.Fatal(err)
    }

    result, err := json.Marshal(raw_result.Value)
    if err != nil {
        log.Fatal(err)
    }

    return string(result)
}


func main() {
    app := martini.Classic()
    app.Use(render.Renderer())

    app.Get("/", func(render render.Render) {
        render.JSON(200, map[string]interface{}{"greeting": "Hello, I'm your API!"})
    })

    app.Post("/api/v1/get", func(render render.Render) {
        render.JSON(200, db_init())
    })

    app.Post("/api/v1/set", func(render render.Render) {
        render.JSON(200, map[string]interface{}{"success": true, "error": nil})
    })

    fmt.Printf("App starting!\n")

    app.Run()
}

