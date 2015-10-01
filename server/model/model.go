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

  "expense_mon/server/helpers"
)


type Main struct {
  Helpers helpers.Main
  MongoCollection *mgo.Collection
  MongoSession *mgo.Session
  DBExpense
  DBExpenseComment
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
  self.MongoCollection = session.DB("test").C("money_mon")

  self.Helpers.LogSimpleMessage("Mongo ready")
}

type DBExpense struct {
  Date time.Time
  Value int
}

type DBExpenseComment struct {
  Date time.Time
  Value int
  Comment string
}

func (self *Main) Get() map[string]interface{} {
  dbExpenses := []DBExpenseComment{}
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
        "week": weekNumber,
        "value": dbExpenses[dbExpenseItr].Value,
        "comment": dbExpenses[dbExpenseItr].Comment,
      }
    } else {
      apiExpense = map[string]interface{}{
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
    "success": apiExpenses,
    "error": nil,
  }
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
        &DBExpenseComment{date, dbExpense.Value, dbExpense.Comment},
      )
    } else {
      self.MongoCollection.Insert(
        &DBExpense{date, dbExpense.Value},
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
  Date int `json: "date"`
}

const (
  OneDayTimestamp = 86400
)

func (self *Main) Remove(res *http.Request) map[string]interface{} {
  dbExpense := ReqRemove{}

  json.Unmarshal(
    []byte(
      self.ProcessReqBody(res),
    ),
    &dbExpense,
  )

  if dbExpense.Date > 0 {
    rawDate :=
      time.Unix(
        int64(dbExpense.Date),
        0,
      )

    date, err :=
      time.Parse(
        "2006-01-02",
        rawDate.Format("2006-01-02"),
      )
    if err != nil {
      self.Helpers.LogWarning(err)
    }

    startDateIntervalTimestamp := int(date.Unix()) - (int(date.Weekday()) - 1) * OneDayTimestamp
    startDateInterval :=
      time.Unix(
        int64(startDateIntervalTimestamp),
        0,
      )

    endDateInterval :=
      time.Unix(
        int64(startDateIntervalTimestamp + OneDayTimestamp * 7),
        0,
      )

    dbExpenses := []DBExpenseComment{}
    self.MongoCollection.Find(nil).All(&dbExpenses)

    // for dbExpenseItr := 0; dbExpenseItr < len(dbExpenses); dbExpenseItr++ {
    //   if dbExpenses[dbExpenseItr].Date.Unix()
    // }

    // self.MongoCollection.Remove()

    // No generics for common method T_T
    fmt.Printf(
      "%s | Removed from DB: %s\n",
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

func (self *Main) ProcessReqBody(res *http.Request) string {
  bodyUInt8, err := ioutil.ReadAll(res.Body)
  if err != nil {
    self.Helpers.LogWarning(err)
  }

  return strings.Replace(string(bodyUInt8), "'", "\"", -1)
}