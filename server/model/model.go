package model

import (
    "fmt"
    // "reflect"
    "log"
    "time"

    "gopkg.in/mgo.v2"
)


type Main struct {
    MongoCollection *mgo.Collection
    MongoSession *mgo.Session
    Expense
}

func (self *Main) Init() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }

    self.MongoSession = session
    self.MongoCollection = session.DB("test").C("money_mon")

    fmt.Printf("Mongo ready\n")
}

type Expense struct {
    Date time.Time
    Value int
    Comment string
}

func (self *Main) Get() map[string]interface{} {
    db_expenses := []Expense{}
    self.MongoCollection.Find(nil).All(&db_expenses)

    api_expenses_month := []map[string]interface{}{}
    api_expenses_year := [][]map[string]interface{}{}
    api_expenses := [][][]map[string]interface{}{}

    current_loop_month := db_expenses[0].Date.Month()
    current_loop_year := db_expenses[0].Date.Year()

    for db_expense_itr := 0; db_expense_itr < len(db_expenses); db_expense_itr++ {
        if db_expenses[db_expense_itr].Date.Month() != current_loop_month {
            api_expenses_year = append(api_expenses_year, api_expenses_month)

            if db_expenses[db_expense_itr].Date.Year() != current_loop_year {
                api_expenses = append(api_expenses, api_expenses_year)

                api_expenses_year = [][]map[string]interface{}{}
            }

            api_expenses_month = []map[string]interface{}{}
        }

        week_number := db_expenses[db_expense_itr].Date.Day() / 7

        if week_number == 0 {
            week_number = 1
        }

        api_expense := map[string]interface{}{}

        if len(db_expenses[db_expense_itr].Comment) > 0 {
            api_expense = map[string]interface{}{
                "week": week_number,
                "value": db_expenses[db_expense_itr].Value,
                "comment": db_expenses[db_expense_itr].Comment,
            }
        } else {
            api_expense = map[string]interface{}{
                "week": week_number,
                "value": db_expenses[db_expense_itr].Value,
            }
        }

        api_expenses_month = append(api_expenses_month, api_expense)

        current_loop_month = db_expenses[db_expense_itr].Date.Month()
        current_loop_year = db_expenses[db_expense_itr].Date.Year()

        // Last iteration
        if db_expense_itr + 1 == len(db_expenses) {
            api_expenses = append(api_expenses, api_expenses_year)
        }
    }

    return map[string]interface{}{
        "success": api_expenses,
        "error": nil,
    }
}

func (self *Main) Set() string {
    // data := self.Data{Value: local_data}

    // err = collection.Insert(data)
    // if err != nil {
    //     log.Fatal(err)
    // }
    return "test"
}