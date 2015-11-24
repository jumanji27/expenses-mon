package expensesModel

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
  Expenses [][][]map[string]interface{}
  GetDBExpense
  SetReq
  SetDBExpenseRequired
  SetDBExpense
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

type GetDBExpense struct {
  Id bson.ObjectId `bson:"_id"`
  Date time.Time
  Value int
  Comment string
  averageUSDRUBRate float32
}

const (
  UnitMeasure = 5000
  Currency = "RUB"
  DayTimestamp = 1 * 24 * 60 * 60
  WeekTimestamp = 7 * DayTimestamp
  MinMonthTimestamp = 28 * DayTimestamp
  MonthsInYear = 12
)

func (self *Main) GetHandler() map[string]interface{} {
  dbExpenses := []GetDBExpense{}
  self.MongoCollection.Find(nil).All(&dbExpenses)

  dbExpensesLength := len(dbExpenses)

  if dbExpensesLength > 0 {
    self.Expenses = [][][]map[string]interface{}{}
    expensesYear := [][]map[string]interface{}{}
    expensesMonth := []map[string]interface{}{}

    apiExpenses := [][][]map[string]interface{}{}
    apiExpensesYear := [][]map[string]interface{}{}
    apiExpensesMonth := []map[string]interface{}{}

    // In first iteration == current
    prevMonth := dbExpenses[0].Date.Month()
    prevYear := dbExpenses[0].Date.Year()

    var weekNumber int
    var prevWeekNumber int
    var firstDayOfMonthIsSunday bool
    var gap int
    var fullYearLoop bool

    // Move to db
    // averageUSDRUBRate :=
    //   map[int]float32{
    //     2013: 31.9,
    //     2014: 38.6,
    //   }

    monthOffset := int(dbExpenses[0].Date.Weekday()) - 1  // Sunday is last weekday in EU

    // Sunday is last weekday in EU
    if monthOffset < 0 {
      monthOffset = 6
      firstDayOfMonthIsSunday = true
    }

    for key, dbExpense := range dbExpenses {
      date := dbExpense.Date
      timestamp := int(date.Unix())
      year := date.Year()
      month := date.Month()
      day := date.Day()

      if month != prevMonth {
        formatedMonths :=
          self.addEmptyExpenses(
            expensesMonth,
            apiExpensesMonth,
            true,
            map[string]int{
              "weekNumber": weekNumber,
              "timestamp": timestamp,
              "gap": gap,
            },
          )

        expensesYear = append(expensesYear, formatedMonths[0])
        apiExpensesYear = append(apiExpensesYear, formatedMonths[1])

        if year != prevYear {
          self.Expenses = append(self.Expenses, expensesYear)
          apiExpenses = append(apiExpenses, apiExpensesYear)

          expensesYear = [][]map[string]interface{}{}
          apiExpensesYear = [][]map[string]interface{}{}

          fullYearLoop = true
        } else {
          fullYearLoop = false
        }

        firstMonthDay :=
          time.Unix(
            int64(timestamp - (day - 1) * DayTimestamp),
            0,
          )

        monthOffset = int(firstMonthDay.Weekday()) - 1

        // Sunday is last weekday in EU
        if monthOffset < 0 {
          monthOffset = 6
          firstDayOfMonthIsSunday = true
        }

        expensesMonth = []map[string]interface{}{}
        apiExpensesMonth = []map[string]interface{}{}
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
      commentLength := len(dbExpense.Comment)

      expense :=
        map[string]interface{}{
          "id": dbExpense.Id,
          "date": dbExpense.Date,
          "value": dbExpense.Value,
          "comment": dbExpense.Comment,
          "year_average_usd_rub_rate": dbExpense.averageUSDRUBRate,
        }

      if dbExpense.averageUSDRUBRate > 0 {
        if commentLength > 0 {
          apiExpense =
            map[string]interface{}{
              "id": dbExpense.Id,
              "value": dbExpense.Value,
              "comment": dbExpense.Comment,
              "year_average_usd_rub_rate": dbExpense.averageUSDRUBRate,
            }
        } else {
          apiExpense =
            map[string]interface{}{
              "id": dbExpense.Id,
              "value": dbExpense.Value,
              "year_average_usd_rub_rate": dbExpense.averageUSDRUBRate,
            }
        }
      } else if commentLength > 0 {
        apiExpense =
          map[string]interface{}{
            "id": dbExpense.Id,
            "value": dbExpense.Value,
            "comment": dbExpense.Comment,
          }
      } else {
        apiExpense =
          map[string]interface{}{
            "id": dbExpense.Id,
            "value": dbExpense.Value,
          }
      }

      gap = weekNumber - prevWeekNumber

      if gap > 1 {
        for extraItr := 1; extraItr < gap; extraItr++ {
          expensesMonth =
            append(
              expensesMonth,
              map[string]interface{}{
                "id": bson.NewObjectId(),
                "date":
                  time.Unix(
                    int64(timestamp - gap * WeekTimestamp),
                    0,
                  ),
              },
            )

          apiExpensesMonth =
            append(
              apiExpensesMonth,
              map[string]interface{}{
                "id": bson.NewObjectId(),
              },
            )
        }
      } else if gap < 0 {
        for extraItr := 1; extraItr < weekNumber; extraItr++ {
          expensesMonth =
            append(
              expensesMonth,
              map[string]interface{}{
                "id": bson.NewObjectId(),
                "date":
                  time.Unix(
                    int64(timestamp - extraItr * WeekTimestamp),
                    0,
                  ),
              },
            )

          apiExpensesMonth =
            append(
              apiExpensesMonth,
              map[string]interface{}{
                "id": bson.NewObjectId(),
              },
            )
        }
      }

      expensesMonth = append(expensesMonth, expense)
      apiExpensesMonth = append(apiExpensesMonth, apiExpense)

      // Non full year
      if key + 1 == dbExpensesLength && fullYearLoop != true {
        formatedMonths :=
          self.addEmptyExpenses(
            expensesMonth,
            apiExpensesMonth,
            false,
            map[string]int{
              "weekNumber": weekNumber,
              "timestamp": timestamp,
              "gap": gap,
            },
          )

        expensesYear = append(expensesYear, formatedMonths[0])
        apiExpensesYear = append(apiExpensesYear, formatedMonths[1])

        now := time.Now()
        timestampGap := int(now.Unix()) - timestamp

        // Fill empty months
        if timestampGap > MinMonthTimestamp {
          gap := timestampGap / MinMonthTimestamp

          if gap > 0 && gap < MonthsInYear {
            for extraItr := 0; extraItr < gap; extraItr++ {
              emptyMonth := []map[string]interface{}{}
              emptyApiMonth := []map[string]interface{}{}

              for extraExpenseItr := 0; extraExpenseItr < WeeksInMonth; extraExpenseItr++ {
                firstDayOfMonthTimestamp := timestamp - (day - 1) * DayTimestamp
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

                emptyMonth =
                  append(
                    emptyMonth,
                    map[string]interface{}{
                      "id": bson.NewObjectId(),
                      "date":
                        time.Unix(
                          int64(timestamp + offset * DayTimestamp + extraExpenseItr * WeekTimestamp),
                          0,
                        ),
                    },
                  )

                emptyApiMonth =
                  append(
                    emptyApiMonth,
                    map[string]interface{}{
                      "id": bson.NewObjectId(),
                    },
                  )
              }

              expensesYear = append(expensesYear, emptyMonth)
              apiExpensesYear = append(apiExpensesYear, emptyApiMonth)
            }
          }
        }

        self.Expenses = append(self.Expenses, expensesYear)
        apiExpenses = append(apiExpenses, apiExpensesYear)

        // We haven't optional behavior for empty months > 12 (empty years). If we need this logic, it'll be here
      }

      prevMonth = month
      prevYear = year
      prevWeekNumber = weekNumber
    }

    apiExpensesInverse := [][][]map[string]interface{}{}

    for key := range apiExpenses {
      apiExpensesInverse =
        append(
          apiExpensesInverse,
          apiExpenses[len(apiExpenses) - key - 1],
        )
    }

    return map[string]interface{}{
      "success":
        map[string]interface{}{
          "expenses": apiExpensesInverse,
          "unit_measure": UnitMeasure,
          "currency": Currency,
        },
    }
  } else {
    status := "DB is empty"

    self.Helpers.CreateEvent("Warning", status)

    return map[string]interface{}{
      "error": status,
    }
  }
}

func (self *Main) addEmptyExpenses(
  month []map[string]interface{}, apiMonth []map[string]interface{}, redefineWeek bool, params map[string]int,
  ) [2][]map[string]interface{} {
    addition := WeeksInMonth - params["weekNumber"]

    if addition > 0 {
      for extraItr := 0; extraItr < addition; extraItr++ {
        month =
          append(
            month,
            map[string]interface{}{
              "id": bson.NewObjectId(),
              "date":
                time.Unix(
                  int64(params["timestamp"] - params["gap"] * WeekTimestamp),
                  0,
                ),
            },
          )

        apiMonth =
          append(
            apiMonth,
            map[string]interface{}{
              "id": bson.NewObjectId(),
            },
          )

        if redefineWeek == true && extraItr + 1 == addition {
          params["weekNumber"] = params["weekNumber"] + extraItr + 1
        }
      }
    }

    result := [2][]map[string]interface{}{}

    result[0] = month
    result[1] = apiMonth

    return result
}

type SetReq struct {
  Action string `json: "action"`
  Id string `json: "id"`
  Comment string `json: "comment"`
}

type SetDBExpenseRequired struct {
  Id bson.ObjectId
  Date time.Time
  Value int
}

type SetDBExpense struct {
  Id bson.ObjectId
  Date time.Time
  Value int
  Comment string
}

func (self *Main) SetHandler(res *http.Request) map[string]interface{} {
  reqExpense := SetReq{}

  json.Unmarshal(
    []byte(
      self.ProcessReqBody(res),
    ),
    &reqExpense,
  )

  if len(reqExpense.Action) > 0 && len(reqExpense.Id) > 0 {
    var value int

    dbExpense := GetDBExpense{}

    self.MongoCollection.Find(
      bson.M{
        "_id": bson.ObjectIdHex(reqExpense.Id),
      },
    ).One(&dbExpense)

    if len(dbExpense.Id) > 0 {
      if reqExpense.Action == "add" {
        value = dbExpense.Value + 1
      } else if reqExpense.Action == "remove" {
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

      if value > 0 {
        self.MongoCollection.Update(
          bson.M{
            "_id": dbExpense.Id,
          },
          bson.M{
            "$set": expense,
          },
        )

        self.Helpers.CreateEvent("Log", "Updated expense")
      } else {
        self.MongoCollection.Remove(
          bson.M{
            "_id": dbExpense.Id,
          },
        )

        self.Helpers.CreateEvent("Log", "Deleted expense")
      }
    } else {
      value = 1

      for _, year := range self.Expenses {
        for _, month := range year {
          for _, expense := range month {
            expense["value"] = value

            expenseId := expense["id"].(bson.ObjectId)

            // WTF DATE?

            if expense["date"] != nil && bson.ObjectId.Hex(expenseId) == reqExpense.Id {
              date := expense["date"].(time.Time)

              if expense["comment"] == nil {
                self.MongoCollection.Insert(
                  bson.M{
                    "_id": expenseId,
                    "date": date,
                    "value": value,
                  },
                )
              } else {
                self.MongoCollection.Insert(
                  bson.M{
                    "_id": expenseId,
                    "date": date,
                    "value": value,
                    "comment": expense["comment"].(string),
                  },
                )
              }
            }
          }
        }
      }

      self.Helpers.CreateEvent("Log", "Added expense")
    }

    return map[string]interface{}{
      "success": true,
    }
  } else {
    self.Helpers.CreateEvent("Log", "Failed request, validation error")

    return map[string]interface{}{
      "error": "Data validation error",
    }
  }
}

func (self *Main) ProcessReqBody(res *http.Request) string {
  bodyUInt8, err := ioutil.ReadAll(res.Body)
  if err != nil {
    self.Helpers.CreateEvent("Warning", err.Error())
  }

  return strings.Replace(
    string(bodyUInt8),
    "'",
    "\"",
    -1,
  )
}