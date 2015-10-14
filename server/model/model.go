package model

import (
  "fmt"
  // "reflect"
  "time"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "strings"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"

  "expense_mon/server/helpers"
)


type Main struct {
  Helpers helpers.Main
  MongoCollection *mgo.Collection
  MongoSession *mgo.Session
  DBExpense
  DBExpenseRequred
  DBExpenseSet
  DBExpenseSetRequred
  ReqSet
}

const (
  LogTimeFormat = "02 Jan 2006 15:04:05"
)


func (self *Main) Init() {
  self.Helpers = helpers.Main{}

  session, err := mgo.Dial("localhost:27017")
  if err != nil {
    self.Helpers.LogError(err)
  }

  self.MongoSession = session
  self.MongoCollection = session.DB("money_mon").C("index")

  self.Helpers.LogSimpleMessage("Mongo ready")
}

type DBExpense struct {
  Id bson.ObjectId `bson:"_id"`
  Date time.Time
  Value int
  Comment string
}

type DBExpenseRequred struct {
  Id bson.ObjectId `bson:"_id"`
  Date time.Time
  Value int
}

const (
  UnitMeasure = 5000
)

func (self *Main) Get() map[string]interface{} {
  dbExpenses := []DBExpense{}
  self.MongoCollection.Find(nil).All(&dbExpenses)

  apiExpensesMonth := []map[string]interface{}{}
  apiExpensesYear := [][]map[string]interface{}{}
  apiExpenses := [][][]map[string]interface{}{}

  currentLoopMonth := dbExpenses[0].Date.Month()
  currentLoopYear := dbExpenses[0].Date.Year()

  var fullYearLoop bool

  // Loop is depended from DB struct (year must begin from january)
  for dbExpenseItr := 0; dbExpenseItr < len(dbExpenses); dbExpenseItr++ {
    if dbExpenses[dbExpenseItr].Date.Month() != currentLoopMonth {
      apiExpensesYear = append(apiExpensesYear, apiExpensesMonth)

      if dbExpenses[dbExpenseItr].Date.Year() != currentLoopYear {
        apiExpenses = append(apiExpenses, apiExpensesYear)

        apiExpensesYear = [][]map[string]interface{}{}
      }

      apiExpensesMonth = []map[string]interface{}{}
      fullYearLoop = true
    }

    weekNumber := dbExpenses[dbExpenseItr].Date.Day() / 7

    if weekNumber == 0 {
      weekNumber = 1
    }

    apiExpense := map[string]interface{}{}

    if len(dbExpenses[dbExpenseItr].Comment) > 0 {
      apiExpense = map[string]interface{}{
        "id": dbExpenses[dbExpenseItr].Id,
        "week": weekNumber,
        "value": dbExpenses[dbExpenseItr].Value,
        "comment": dbExpenses[dbExpenseItr].Comment,
      }
    } else {
      apiExpense = map[string]interface{}{
        "id": dbExpenses[dbExpenseItr].Id,
        "week": weekNumber,
        "value": dbExpenses[dbExpenseItr].Value,
      }
    }

    apiExpensesMonth = append(apiExpensesMonth, apiExpense)

    currentLoopMonth = dbExpenses[dbExpenseItr].Date.Month()
    currentLoopYear = dbExpenses[dbExpenseItr].Date.Year()

    // Last iteration
    if dbExpenseItr + 1 == len(dbExpenses) {
      if fullYearLoop != true {
        apiExpensesYear = append(apiExpensesYear, apiExpensesMonth)
      }

      apiExpenses = append(apiExpenses, apiExpensesYear)

      // For first empty month new year
      // Empty array instead null in response — bad API design T_T
      if currentLoopYear != time.Now().Year() {
        apiExpenses = append(apiExpenses, [][]map[string]interface{}{})
      }
    }
  }

  return map[string]interface{}{
    "success":
      map[string]interface{}{
        "unit_measure": UnitMeasure,
        "expenses": apiExpenses,
      },
    "error": nil,
  }
}

type DBExpenseSet struct {
  Date time.Time
  Value int
  Comment string
}

type DBExpenseSetRequred struct {
  Date time.Time
  Value int
}

type ReqSet struct {
  Date int `json: "date"`
  Value int `json: "value"`
  Comment string `json: "comment"`
}

func (self *Main) Set(res *http.Request) map[string]interface{} {
  dbExpense := ReqSet{}

  json.Unmarshal(
    []byte(
      self.ProcessReqBody(res),
    ),
    &dbExpense,
  )

  if dbExpense.Value > 0 {
    date64 := int64(dbExpense.Date)
    date := time.Unix(date64, 0)

    if len(dbExpense.Comment) > 0 {
      self.MongoCollection.Insert(
        &DBExpenseSet{date, dbExpense.Value, dbExpense.Comment},
      )
    } else {
      self.MongoCollection.Insert(
        &DBExpenseSetRequred{date, dbExpense.Value},
      )
    }

    // No generics for common method T_T
    fmt.Printf(
      "%s | Added to DB: %s\n",
      time.Now().Format(LogTimeFormat),
      dbExpense,
    )

    return map[string]interface{}{
      "success": true,
      "error": nil,
    }
  } else {
    self.Helpers.LogSimpleMessage("Failed request, validation error")

    return map[string]interface{}{
      "success": nil,
      "error": "Data validation error",
    }
  }
}

type ReqRemove struct {
  Id string `json: "id"`
}

const (
  OneDayTimestamp = 86400
)

func (self *Main) Remove(res *http.Request) map[string]interface{} {
  reqExpense := ReqRemove{}

  json.Unmarshal(
    []byte(
      self.ProcessReqBody(res),
    ),
    &reqExpense,
  )

  if len(reqExpense.Id) > 0 {
    dbExpense := DBExpense{}
    self.MongoCollection.Find(
      bson.M{
        "_id": bson.ObjectIdHex(reqExpense.Id),
      },
    ).One(&dbExpense)

    self.MongoCollection.Remove(
      bson.M{
        "_id": bson.ObjectIdHex(reqExpense.Id),
      },
    )

    if len(dbExpense.Id) > 0 {
      // No generics for common method T_T
      fmt.Printf(
        "%s | Removed from DB: %s\n",
        time.Now().Format(LogTimeFormat),
        reqExpense,
      )

      return map[string]interface{}{
        "success": true,
        "error": nil,
      }
    } else {
      return map[string]interface{}{
        "success": nil,
        "error": "Did't found this expense",
      }
    }
  } else {
    self.Helpers.LogSimpleMessage("Failed request, validation error")

    return map[string]interface{}{
      "success": nil,
      "error": "Data validation error",
    }
  }
}

func (self *Main) ProcessReqBody(res *http.Request) string {
  bodyUInt8, err := ioutil.ReadAll(res.Body)
  if err != nil {
    self.Helpers.LogWarning(err)
  }

  return strings.Replace(string(bodyUInt8), "'", "\"", -1)
}