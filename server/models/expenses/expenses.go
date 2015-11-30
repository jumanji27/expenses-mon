package expensesModel

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
  Expenses [][][]map[string]interface{}
  APIExpenses [][][]map[string]interface{}
  GetDBExpense
  SetReq
  SetDBExpenseRequired
  SetDBExpense
}

const (
  WeeksInMonth = 5 // UI restrictions
  DayTimestamp = 1 * 24 * 60 * 60
  WeekTimestamp = 7 * DayTimestamp
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

  self.formExpenses()
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
  EmptyDBErrorMessage = "DB is empty"
)

func (self *Main) GetHandler() map[string]interface{} {
  if len(self.APIExpenses) > 0 {
    inverseExpenses := [][][]map[string]interface{}{}

    for key := range self.APIExpenses {
      inverseExpenses =
        append(
          inverseExpenses,
          self.APIExpenses[len(self.APIExpenses) - key - 1],
        )
    }

    return map[string]interface{}{
      "success":
        map[string]interface{}{
          "expenses": inverseExpenses,
          "unit_measure": UnitMeasure,
          "currency": Currency,
        },
    }
  } else {
    self.Helpers.CreateEvent("Warning", EmptyDBErrorMessage)

    return map[string]interface{}{
      "error": EmptyDBErrorMessage,
    }
  }
}

const (
  MonthsInYear = 12
  MinMonthTimestamp = 28 * DayTimestamp
)

func (self *Main) formExpenses() {
  dbExpenses := []GetDBExpense{}
  self.MongoCollection.Find(nil).Sort("date").All(&dbExpenses)

  dbExpensesLength := len(dbExpenses)

  if dbExpensesLength > 0 {
    expensesYear := [][]map[string]interface{}{}
    expensesMonth := []map[string]interface{}{}

    APIExpensesYear := [][]map[string]interface{}{}
    APIExpensesMonth := []map[string]interface{}{}

    // In first iteration == current
    prevDate := dbExpenses[0].Date

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

    if len(self.Expenses) > 0 {
      self.Expenses = [][][]map[string]interface{}{}
    }

    if len(self.APIExpenses) > 0 {
      self.APIExpenses = [][][]map[string]interface{}{}
    }

    for key, dbExpense := range dbExpenses {
      var emptyId bson.ObjectId

      date := dbExpense.Date
      timestamp := int(date.Unix())
      year := date.Year()
      month := date.Month()
      day := date.Day()

      if month != prevDate.Month() {
        formatedMonths := self.addEmptyExpenses(expensesMonth, APIExpensesMonth, prevDate, true, weekNumber)

        expensesYear = append(expensesYear, formatedMonths[0])
        APIExpensesYear = append(APIExpensesYear, formatedMonths[1])

        if year != prevDate.Year() {
          self.Expenses = append(self.Expenses, expensesYear)
          self.APIExpenses = append(self.APIExpenses, APIExpensesYear)

          expensesYear = [][]map[string]interface{}{}
          APIExpensesYear = [][]map[string]interface{}{}

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
        APIExpensesMonth = []map[string]interface{}{}
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

      APIExpense := map[string]interface{}{}
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
          APIExpense =
            map[string]interface{}{
              "id": dbExpense.Id,
              "value": dbExpense.Value,
              "comment": dbExpense.Comment,
              "year_average_usd_rub_rate": dbExpense.averageUSDRUBRate,
            }
        } else {
          APIExpense =
            map[string]interface{}{
              "id": dbExpense.Id,
              "value": dbExpense.Value,
              "year_average_usd_rub_rate": dbExpense.averageUSDRUBRate,
            }
        }
      } else if commentLength > 0 {
        APIExpense =
          map[string]interface{}{
            "id": dbExpense.Id,
            "value": dbExpense.Value,
            "comment": dbExpense.Comment,
          }
      } else {
        APIExpense =
          map[string]interface{}{
            "id": dbExpense.Id,
            "value": dbExpense.Value,
          }
      }

      emptyId = bson.NewObjectId()
      gap = weekNumber - prevWeekNumber

      if gap > 1 {
        for itr := 1; itr < gap; itr++ {
          emptyDate :=
            time.Unix(
              int64(
                timestamp - (gap - itr) * WeekTimestamp,
              ),
              0,
            )

          fmt.Println("1")
          fmt.Println(emptyDate)

          expensesMonth =
            append(
              expensesMonth,
              map[string]interface{}{
                "id": emptyId,
                "date": emptyDate,
              },
            )

          APIExpensesMonth =
            append(
              APIExpensesMonth,
              map[string]interface{}{
                "id": emptyId,
              },
            )
        }
      } else if gap < 0 {
        for itr := 1; itr < weekNumber; itr++ {
          emptyDate :=
            time.Unix(
              int64(
                timestamp - (weekNumber - itr) * WeekTimestamp,
              ),
              0,
            )

          if emptyDate.Month() != date.Month() {
            day := int(date.Day())

            emptyDate =
              time.Unix(
                int64(
                  timestamp - (day - 1) * DayTimestamp,
                ),
                0,
              )
          }

          fmt.Println("2")
          fmt.Println(emptyDate)

          expensesMonth =
            append(
              expensesMonth,
              map[string]interface{}{
                "id": emptyId,
                "date": emptyDate,
              },
            )

          APIExpensesMonth =
            append(
              APIExpensesMonth,
              map[string]interface{}{
                "id": emptyId,
              },
            )
        }
      }

      expensesMonth = append(expensesMonth, expense)
      APIExpensesMonth = append(APIExpensesMonth, APIExpense)

      fmt.Println("DB")
      fmt.Println(expense["date"])

      // Non full year
      if key + 1 == dbExpensesLength && fullYearLoop != true {
        formatedMonths := self.addEmptyExpenses(expensesMonth, APIExpensesMonth, prevDate, false, weekNumber)

        expensesYear = append(expensesYear, formatedMonths[0])
        APIExpensesYear = append(APIExpensesYear, formatedMonths[1])

        now := time.Now()
        timestampGap := int(now.Unix()) - timestamp

        // Fill empty months
        if timestampGap > MinMonthTimestamp {
          gap := timestampGap / MinMonthTimestamp

          if gap > 0 && gap < MonthsInYear {
            for itr := 0; itr < gap; itr++ {
              emptyMonth := []map[string]interface{}{}
              emptyApiMonth := []map[string]interface{}{}

              for extraItr := 0; extraItr < WeeksInMonth; extraItr++ {
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
                      "id": emptyId,
                      "date":
                        time.Unix(
                          int64(timestamp + offset * DayTimestamp + extraItr * WeekTimestamp),
                          0,
                        ),
                    },
                  )

                emptyApiMonth =
                  append(
                    emptyApiMonth,
                    map[string]interface{}{
                      "id": emptyId,
                    },
                  )
              }

              expensesYear = append(expensesYear, emptyMonth)
              APIExpensesYear = append(APIExpensesYear, emptyApiMonth)
            }
          }
        }

        self.Expenses = append(self.Expenses, expensesYear)
        self.APIExpenses = append(self.APIExpenses, APIExpensesYear)

        // We haven't optional behavior for empty months > 12 (empty years). If we need this logic, it'll be here
      }

      prevDate = date
      prevWeekNumber = weekNumber
    }
  }
}

func (self *Main) addEmptyExpenses(
  month []map[string]interface{}, APIMonth []map[string]interface{},
  prevDate time.Time, redefineWeek bool, weekNumber int,
  ) [2][]map[string]interface{} {
    id := bson.NewObjectId()
    addition := WeeksInMonth - weekNumber

    if addition > 0 {
      for itr := 0; itr < addition; itr++ {
        date :=
          time.Unix(
            int64(
              int(prevDate.Unix()) + (addition - itr) * WeekTimestamp,
            ),
            0,
          )

        fmt.Println("4")
        fmt.Println(date)

        month =
          append(
            month,
            map[string]interface{}{
              "id": id,
              "date": date,
            },
          )

        APIMonth =
          append(
            APIMonth,
            map[string]interface{}{
              "id": id,
            },
          )

        if redefineWeek == true && itr + 1 == addition {
          weekNumber = weekNumber + itr + 1
        }
      }
    }

    result := [2][]map[string]interface{}{}

    result[0] = month
    result[1] = APIMonth

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
      // Next receive any value instead of action
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

      return map[string]interface{}{
        "success": true,
      }
    } else {
      var matchExpense map[string]interface{}

      value = 1

      for key, year := range self.Expenses {
        for monthKey, month := range year {
          for expenseKey, expense := range month {
            expenseId := expense["id"].(bson.ObjectId)

            if bson.ObjectId.Hex(expenseId) == reqExpense.Id {
              expense["value"] = value
              self.APIExpenses[key][monthKey][expenseKey]["value"] = value

              matchExpense = expense
              break
            }
          }
        }
      }

      expenseId := matchExpense["id"].(bson.ObjectId)

      if len(bson.ObjectId.Hex(expenseId)) > 0 {
        date := matchExpense["date"].(time.Time)

        if matchExpense["comment"] == nil {
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
              "comment": matchExpense["comment"].(string),
            },
          )
        }

        self.Helpers.CreateEvent("Log", "Added expense")

        return map[string]interface{}{
          "success": true,
        }
      } else {
        return map[string]interface{}{
          "error": "Not found this expense ID",
        }
      }
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