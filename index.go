package main


import (
    "fmt"
    "log"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
)


type Value struct {
    key string
    value int
}


func db_init() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        panic(err)
    }

    defer session.Close()

    c := session.DB("test").C("test")


    const key string = "value"
    var data [1][3][12][5]Value

    for i := 0; i < 3; i++ {
        for j := 0; j < 12; j++ {
            for k := 0; k < 5; k++ {
                data[0][i][j][k] = Value{key, k}
            }
        }
    }

    err = c.Insert(data)
    if err != nil {
        log.Fatal(err)
    }
}

// func db_find() {
//     result := Data{}
//     err = c.Find(bson.M{"name": "Ale"}).One(&result)
//     if err != nil {
//         log.Fatal(err)
//     }

//     fmt.Println("Phone:", result.Phone)
// }


func main() {
    app := martini.Classic()
    app.Use(render.Renderer())

    db_init()

    // app.Post("/api/v1/get", func(render render.Render) {
    //     render.JSON(200, db_find())
    // })

    app.Post("/api/v1/set", func(render render.Render) {
        render.JSON(200, map[string]interface{}{"success": true, "error": nil})
    })

    fmt.Printf("App starting!\n")

    app.Run()
}

