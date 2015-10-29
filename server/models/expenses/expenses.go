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

  "expenses-mon/server/helpers"
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
  OneDayTimestamp = 86400
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
  Currency = "RUB"
)

func (self *Main) Get() map[string]interface{} {
  dbExpenses := []DBExpense{}
  self.MongoCollection.Find(nil).All(&dbExpenses)

  dbExpensesLength := len(dbExpenses)

  apiExpensesMonth := []map[string]interface{}{}
  apiExpensesYear := [][]map[string]interface{}{}
  apiExpenses := [][][]map[string]interface{}{}

  currentLoopMonth := dbExpenses[0].Date.Month()
  currentLoopYear := dbExpenses[0].Date.Year()
  monthOffset := int(dbExpenses[0].Date.Weekday()) - 1  // Sunday is last weekday in EU

  var fullYearLoop bool
  var firstDayOfMonthIsSunday bool

  // Handle fill map
  averageUSDRUBRate :=
    map[int]float32{
      2013: 31.9,
      2014: 38.6,
    }

  // Sunday is last weekday in EU
  if monthOffset < 0 {
    monthOffset = 6
    firstDayOfMonthIsSunday = true
  }

  // Loop is depended from DB struct (year must begin from january)
  for dbExpenseItr := 0; dbExpenseItr < dbExpensesLength; dbExpenseItr++ {
    month := dbExpenses[dbExpenseItr].Date.Month()
    year := dbExpenses[dbExpenseItr].Date.Year()

    if month != currentLoopMonth {
      apiExpensesYear = append(apiExpensesYear, apiExpensesMonth)

      if year != currentLoopYear {
        apiExpenses = append(apiExpenses, apiExpensesYear)

        apiExpensesYear = [][]map[string]interface{}{}
      }

      firstMonthDay :=
        time.Unix(
          int64(
            int(dbExpenses[dbExpenseItr].Date.Unix()) - (dbExpenses[dbExpenseItr].Date.Day() - 1) * OneDayTimestamp,
          ),
          0,
        )

      monthOffset = int(firstMonthDay.Weekday()) - int(firstMonthDay.Day())

      // Sunday is last weekday in EU
      if monthOffset < 0 {
        monthOffset = 6
        firstDayOfMonthIsSunday = true
      }

      apiExpensesMonth = []map[string]interface{}{}
      fullYearLoop = true
    }

    day := dbExpenses[dbExpenseItr].Date.Day()
    var weekNumber int

    if firstDayOfMonthIsSunday == true && day == 1 {
      weekNumber = 1
    } else {
      weekNumber = (monthOffset + day) / 7 + 1
    }

    // UI possible restrictions
    if weekNumber > 5 {
      weekNumber = 5
    }

    firstDayOfMonthIsSunday = false

    apiExpense := map[string]interface{}{}
    id := dbExpenses[dbExpenseItr].Id
    comment := dbExpenses[dbExpenseItr].Comment
    value := dbExpenses[dbExpenseItr].Value
    monthInt := int(month)
    commentLength := len(comment)

    if monthInt == 1 && averageUSDRUBRate[year] > 0 && commentLength > 0 {
      apiExpense = map[string]interface{}{
        "id": id,
        "week": weekNumber,
        "value": value,
        "comment": comment,
        "year_average_usd_rub_rate": averageUSDRUBRate[year],
      }
    } else if monthInt == 1 && averageUSDRUBRate[year] > 0 {
      apiExpense = map[string]interface{}{
        "id": id,
        "week": weekNumber,
        "value": value,
        "year_average_usd_rub_rate": averageUSDRUBRate[year],
      }
    } else if commentLength > 0 {
      apiExpense = map[string]interface{}{
        "id": id,
        "week": weekNumber,
        "value": value,
        "comment": comment,
      }
    } else {
      apiExpense = map[string]interface{}{
        "id": id,
        "week": weekNumber,
        "value": value,
      }
    }

    apiExpensesMonth = append(apiExpensesMonth, apiExpense)

    currentLoopMonth = month
    currentLoopYear = year

    // Last iteration
    if dbExpenseItr + 1 == dbExpensesLength {
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
        "expenses": apiExpenses,
        "unit_measure": UnitMeasure,
        "currency": Currency,
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