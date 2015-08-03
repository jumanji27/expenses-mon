package main


import (
    "fmt"
    "log"
    "time"
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
    DBExpense
    API
    APIWeek
}

type DBExpense struct {
    Date time.Time
    Value int
    Comment string
}

type API struct {
    Value [1][][12][5]APIWeek
}

type APIWeek struct {
    Week int
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

    db_expenses := []DBExpense{}
    self.MongoCollection.Find(nil).All(&db_expenses)

    api_expenses_years := make([][]APIWeek, 2) // ???
    api_expenses_weeks := []APIWeek{}          // ???

    year := db_expenses[0].Date.Year()
    year_itr := 0

    for db_expense_itr := 0; db_expense_itr < len(db_expenses); db_expense_itr++ {
        if db_expense_itr == 0 {
            api_expenses_years[db_expense_itr] = api_expenses_weeks
        } else if db_expenses[db_expense_itr].Date.Year() != year {
            year = db_expenses[db_expense_itr].Date.Year()
            year_itr++
            api_expenses_years[year_itr] = api_expenses_weeks
        }
    }

    // REFACTOR AND DEBUG
    for year_itr := 0; year_itr < len(api_expenses_years); year_itr++ {
        fmt.Println(year_itr)

        // for month_itr := 1; month_itr < 12; month_itr++ {
        //     for db_expense_itr := 0; db_expense_itr < len(db_expenses); db_expense_itr++ {
        //         fmt.Println(api_expenses)

        //         // if db_expenses[db_expense_itr].Date.Year() == api_expenses[year_itr] && int(db_expenses[db_expense_itr].Date.Month()) == month_itr && db_expenses[db_expense_itr].Comment {
        //         //     api_expenses[year_itr][month_itr] = APIWeek{1 ,db_expenses[db_expense_itr].Value, db_expenses[db_expense_itr].Comment}
        //         // } else if db_expenses[db_expense_itr].Date.Year() == api_expenses[year_itr] && int(db_expenses[db_expense_itr].Date.Month()) == month_itr {
        //         //     api_expenses[year_itr][month_itr] = APIWeek{1 ,db_expenses[db_expense_itr].Value}
        //         // }
        //     }
        // }
    }


    api_result, err := json.Marshal(db_expenses)
    if err != nil {
        log.Fatal(err)
    }

    return string(api_result)
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
