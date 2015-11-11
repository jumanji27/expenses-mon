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
  ChangeReq
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
  WeekTimestamp = 7 * 24 * 60 * 60
  WeeksInMonth = 5 // UI restrictions
)

func (self *Main) Get() map[string]interface{} {
  dbExpenses := []DBExpense{}
  self.MongoCollection.Find(nil).All(&dbExpenses)

  dbExpensesLength := len(dbExpenses)

  expensesMonth := []map[string]interface{}{}
  expensesYear := [][]map[string]interface{}{}
  expenses := [][][]map[string]interface{}{}


  // In first iteration == current
  prevMonth := dbExpenses[0].Date.Month()
  prevYear := dbExpenses[0].Date.Year()
  prevWeekNumber := 0
  weekNumber := 0
  gap := 0

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

  for dbExpenseItr := 0; dbExpenseItr < dbExpensesLength; dbExpenseItr++ {
    timestamp := int(dbExpenses[dbExpenseItr].Date.Unix())
    year := dbExpenses[dbExpenseItr].Date.Year()
    month := dbExpenses[dbExpenseItr].Date.Month()
    day := dbExpenses[dbExpenseItr].Date.Day()

    if month != prevMonth {
      extraEndOfWeekExpenses := WeeksInMonth - prevWeekNumber

      for extraItr := 0; extraItr < extraEndOfWeekExpenses; extraItr++ {
        week := prevWeekNumber + extraItr + 1

        expensesMonth =
          append(
            expensesMonth,
            map[string]interface{}{
              "week": week,
              "date": timestamp - gap * WeekTimestamp,
            },
          )

        if extraItr + 1 == extraEndOfWeekExpenses {
          prevWeekNumber = week
        }
      }

      expensesYear = append(expensesYear, expensesMonth)

      if year != prevYear {
        expenses = append(expenses, expensesYear)

        expensesYear = [][]map[string]interface{}{}
        fullYearLoop = true
      } else {
        fullYearLoop = false
      }

      firstMonthDay :=
        time.Unix(
          int64(timestamp - (day - 1) * OneDayTimestamp),
          0,
        )

      monthOffset = int(firstMonthDay.Weekday()) - int(firstMonthDay.Day())

      // Sunday is last weekday in EU
      if monthOffset < 0 {
        monthOffset = 6
        firstDayOfMonthIsSunday = true
      }

      expensesMonth = []map[string]interface{}{}
    }

    if firstDayOfMonthIsSunday == true && day == 1 {
      weekNumber = 1
    } else {
      weekNumber = (monthOffset + day) / 7 + 1
    }

    if weekNumber > WeeksInMonth {
      weekNumber = WeeksInMonth
    }

    firstDayOfMonthIsSunday = false

    apiExpense := map[string]interface{}{}
    id := dbExpenses[dbExpenseItr].Id
    comment := dbExpenses[dbExpenseItr].Comment
    value := dbExpenses[dbExpenseItr].Value
    monthInt := int(month)
    monthLength := len(expensesMonth)
    commentLength := len(comment)

    if monthInt == 1 && monthLength == 0 && averageUSDRUBRate[year] > 0 {
      if commentLength > 0 {
        apiExpense =
          map[string]interface{}{
            "id": id,
            "week": weekNumber,
            "value": value,
            "comment": comment,
            "year_average_usd_rub_rate": averageUSDRUBRate[year],
          }
      } else {
        apiExpense =
          map[string]interface{}{
            "id": id,
            "week": weekNumber,
            "value": value,
            "year_average_usd_rub_rate": averageUSDRUBRate[year],
          }
      }
    } else if commentLength > 0 {
      apiExpense =
        map[string]interface{}{
          "id": id,
          "week": weekNumber,
          "value": value,
          "comment": comment,
        }
    } else {
      apiExpense =
        map[string]interface{}{
          "id": id,
          "week": weekNumber,
          "value": value,
        }
    }

    gap = weekNumber - prevWeekNumber

    if gap > 1 {
      for extraItr := 1; extraItr < gap; extraItr++ {
        expensesMonth =
          append(
            expensesMonth,
            map[string]interface{}{
              "week": weekNumber - gap + 1,
              "date": timestamp - gap * WeekTimestamp,
            },
          )
      }
    } else if gap < 0 {
      for extraItr := 1; extraItr < weekNumber; extraItr++ {
        expensesMonth =
          append(
            expensesMonth,
            map[string]interface{}{
              "week": extraItr,
              "date": timestamp - extraItr * WeekTimestamp,
            },
          )
      }
    }

    expensesMonth = append(expensesMonth, apiExpense)

    // Non full year
    if dbExpenseItr + 1 == dbExpensesLength && fullYearLoop != true {
      expensesYear = append(expensesYear, expensesMonth)

      // Fill empty months

      expenses = append(expenses, expensesYear)

      // Fill empty years
    }

    prevMonth = month
    prevYear = year
    prevWeekNumber = weekNumber
  }

  apiExpenses := [][][]map[string]interface{}{}

  for key := range expenses {
    apiExpenses =
      append(
        apiExpenses,
        expenses[len(expenses) - key - 1],
      )
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

func (self *Main) Set(res *http.Request) map[string]interface{} {
  return self.ChangeRecord(res, "set")
}

func (self *Main) Remove(res *http.Request) map[string]interface{} {
  return self.ChangeRecord(res, "remove")
}

func (self *Main) ProcessReqBody(res *http.Request) string {
  bodyUInt8, err := ioutil.ReadAll(res.Body)
  if err != nil {
    self.Helpers.LogWarning(err)
  }

  return strings.Replace(string(bodyUInt8), "'", "\"", -1)
}

type ChangeReq struct {
  Id string `json: "id"`
  Comment string `json: "comment"`
}

func (self *Main) ChangeRecord(res *http.Request, action string) map[string]interface{} {
  reqExpense := ChangeReq{}

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

    if len(dbExpense.Id) > 0 {
      var value int
      var logMessage string

      if action == "set" {
        value = dbExpense.Value + 1
        logMessage = "%s | Added to DB: %s\n"
      } else if action == "remove" {
        value = dbExpense.Value - 1
        logMessage = "%s | Removed from DB: %s\n"
      }

      if len(reqExpense.Comment) > 0 {
        self.MongoCollection.Update(
          bson.M{
            "_id": bson.ObjectIdHex(reqExpense.Id),
          },
          bson.M{
            "$set":
              bson.M{
                "value": value,
                "commit": reqExpense.Comment,
              },
          },
        )
      } else {
        self.MongoCollection.Update(
          bson.M{
            "_id": bson.ObjectIdHex(reqExpense.Id),
          },
          bson.M{
            "$set":
              bson.M{
                "value": value,
              },
          },
        )
      }

      fmt.Printf(
        logMessage,
        time.Now().Format(LogTimeFormat),
        dbExpense,
      )
    } else {
      return map[string]interface{}{
        "success": nil,
        "error": "Did't found this expense",
      }
    }

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