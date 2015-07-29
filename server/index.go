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
    MongoSession *mgo.Session
    DBRecord
    API
    APIWeek
}

type DBRecord struct {
    Date string
    Value int
    Comment string
}

type API struct {
    Value [1][3][12][5]APIWeek
}

type APIWeek struct {
    Value int
    Comment string
}


func (self *Index) db_init() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }

    self.MongoSession = session

    self.MongoCollection = session.DB("test").C("money_mon")
}

func (self *Index) db_get() string {
    defer self.MongoSession.Close()

    db_record := []DBRecord{}
    self.MongoCollection.Find(nil).All(&db_record)



    json_result, err := json.Marshal(db_record)
    if err != nil {
        log.Fatal(err)
    }

    return string(json_result)
}

func(self *Index) db_set() string {
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


func (self *Index) route(app *martini.ClassicMartini) {
    app.Get("/", func(render render.Render) {
        render.JSON(200, map[string]interface{}{
            "success": map[string]interface{}{"greeting": "Hello, I'm your API!"},
            "error": nil,
        })
    })

    app.Post("/api/v1/get", func(render render.Render) {
        render.JSON(200, self.db_get())
    })

    app.Post("/api/v1/set", func(render render.Render) {
        render.JSON(200, self.db_set())
    })
}


func main() {
    martini_app := martini.Classic()
    martini_app.Use(render.Renderer())

    app := Index{}
    app.db_init()
    app.route(martini_app)

    fmt.Printf("App starting!\n")

    martini_app.Run()
}
