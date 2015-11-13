package model

import (
  // "fmt"
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
  WeeksInMonth = 5 // UI restrictions
)


func (self *Main) Init() {
  self.Helpers = helpers.Main{}

  session, err := mgo.Dial("localhost:27017")
  if err != nil {
    self.Helpers.CreateEvent("Error", err.Error())
  }

  self.MongoSession = session
  self.MongoCollection = session.DB("money_mon").C("index")

  self.Helpers.CreateEvent("Log", "Mongo ready")
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
  DayTimestamp = 1 * 24 * 60 * 60
  WeekTimestamp = 7 * DayTimestamp
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

  var prevWeekNumber int
  var gap int
  var fullYearLoop bool

  // Handle fill map
  averageUSDRUBRate :=
    map[int]float32{
      2013: 31.9,
      2014: 38.6,
    }

  for dbExpenseItr := 0; dbExpenseItr < dbExpensesLength; dbExpenseItr++ {
    expense := dbExpenses[dbExpenseItr]
    date := expense.Date
    timestamp := int(date.Unix())
    year := date.Year()
    month := date.Month()

    if month != prevMonth {
      // What if we have month gap one-two and more months?
      extraEndOfWeekExpenses := WeeksInMonth - prevWeekNumber

      for extraItr := 0; extraItr < extraEndOfWeekExpenses; extraItr++ {
        expensesMonth =
          append(
            expensesMonth,
            map[string]interface{}{
              "date": timestamp - gap * WeekTimestamp,
            },
          )

        if extraItr + 1 == extraEndOfWeekExpenses {
          prevWeekNumber = prevWeekNumber + extraItr + 1
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

      expensesMonth = []map[string]interface{}{}
    }

    apiExpense := map[string]interface{}{}
    commentLength := len(expense.Comment)

    if int(month) == 1 && len(expensesMonth) == 0 && averageUSDRUBRate[year] > 0 {
      if commentLength > 0 {
        apiExpense =
          map[string]interface{}{
            "id": expense.Id,
            "value": expense.Value,
            "comment": expense.Comment,
            "year_average_usd_rub_rate": averageUSDRUBRate[year],
          }
      } else {
        apiExpense =
          map[string]interface{}{
            "id": expense.Id,
            "value": expense.Value,
            "year_average_usd_rub_rate": averageUSDRUBRate[year],
          }
      }
    } else if commentLength > 0 {
      apiExpense =
        map[string]interface{}{
          "id": expense.Id,
          "value": expense.Value,
          "comment": expense.Comment,
        }
    } else {
      apiExpense =
        map[string]interface{}{
          "id": expense.Id,
          "value": expense.Value,
        }
    }

    weekNumber := self.GetWeekNumberFromDate(date)

    if prevWeekNumber > 0 {
      gap = weekNumber - prevWeekNumber
    } else {
      gap = weekNumber
    }

    if gap > 1 {
      for extraItr := 1; extraItr < gap; extraItr++ {
        expensesMonth =
          append(
            expensesMonth,
            map[string]interface{}{
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
              "date": timestamp - extraItr * WeekTimestamp,
            },
          )
      }
    }

    expensesMonth = append(expensesMonth, apiExpense)

    // Non full year
    if dbExpenseItr + 1 == dbExpensesLength && fullYearLoop != true {
      expensesYear = append(expensesYear, expensesMonth)

      now := time.Now()

      // Fill empty months
      if now.Year() == year {
        if now.Month() != month {
          gap := int(now.Month()) - int(month)

          if gap > 0 {
            for extraItr := 0; extraItr < gap; extraItr++ {
              month := []map[string]interface{}{}

              for extraExpenseItr := 0; extraExpenseItr < WeeksInMonth; extraExpenseItr++ {
                firstDayOfMonthTimestamp := int(date.Unix()) - (date.Day() - 1) * DayTimestamp
                firstDayOfMonth :=
                  time.Unix(
                    int64(firstDayOfMonthTimestamp),
                    0,
                  )

                firstDayOfMonth.AddDate(0, 1, 0)

                offset := int(firstDayOfMonth.Weekday()) - 1 // Sunday is last weekday in EU

                // Sunday is last weekday in EU
                if offset < 0 {
                  offset = 6
                }

                month =
                  append(
                    month,
                    map[string]interface{}{
                      "date": int(date.Unix()) + offset * DayTimestamp + extraExpenseItr * WeekTimestamp,
                    },
                  )
              }

              expensesYear = append(expensesYear, month)
            }
          }
        }
      }

      expenses = append(expenses, expensesYear)

      // Fill empty years
      if now.Year() != year {
        // ???
      }
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

func (self *Main) GetWeekNumberFromDate(date time.Time) int {
  offset := int(date.Weekday()) - 1 // Sunday is last weekday in EU
  day := date.Day()

  var firstDayOfMonthIsSunday bool
  var weekNumber int

  // Sunday is last weekday in EU
  if offset < 0 {
    offset = 6
    firstDayOfMonthIsSunday = true
  }

  if firstDayOfMonthIsSunday == true && day == 1 {
    weekNumber = 1
  } else {
    weekNumber = (offset + day) / 7 + 1
  }

  if weekNumber > WeeksInMonth {
    weekNumber = WeeksInMonth
  }

  return weekNumber
}

func (self *Main) Set(res *http.Request) map[string]interface{} {
  return self.ChangeRecord(res, "Set")
}

func (self *Main) Remove(res *http.Request) map[string]interface{} {
  return self.ChangeRecord(res, "Remove")
}

func (self *Main) ProcessReqBody(res *http.Request) string {
  bodyUInt8, err := ioutil.ReadAll(res.Body)
  if err != nil {
    self.Helpers.CreateEvent("Warning", err.Error())
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

    // Add and Remove not only Update

    if len(dbExpense.Id) > 0 {
      var value int

      if action == "Set" {
        value = dbExpense.Value + 1
      } else if action == "Remove" {
        value = dbExpense.Value - 1
      }

      var expense bson.M

      if len(reqExpense.Comment) > 0 {
        expense =
          bson.M{
            "value": value,
            "commit": reqExpense.Comment,
          }
      } else {
        expense =
          bson.M{
            "value": value,
          }
      }

      self.MongoCollection.Update(
        bson.M{
          "_id": bson.ObjectIdHex(reqExpense.Id),
        },
        bson.M{
          "$set": expense,
        },
      )

      self.Helpers.CreateEvent("Log", "Updated expense")
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
    self.Helpers.CreateEvent("Log", "Failed request, validation error")

    return map[string]interface{}{
      "success": nil,
      "error": "Data validation error",
    }
  }
}