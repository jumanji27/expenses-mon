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
  APIExpenses [][][]map[string]interface{}
  DBExpense
  SetReq
}

const (
  WeeksInMonth = 5 // UI restrictions
  DayTimestamp = 1 * 24 * 60 * 60
  WeekTimestamp = 7 * DayTimestamp
  DaysInWeek = 7
)


func (self *Main) Init() {
  self.Helpers = helpers.Main{}

  session, err := mgo.Dial("localhost:27017")
  if err != nil {
    self.Helpers.CreateEvent("Error", err.Error())
  }

  self.MongoSession = session
  self.MongoCollection = session.DB("expenses_mon").C("index")

  self.Helpers.CreateEvent("Log", "Mongo ready")
}

type DBExpense struct {
  Id bson.ObjectId `bson:"_id"`
  Date time.Time
  Value int
  Comment string
  YearAverageUSDRUBRate float64 `bson:"year_average_usd_rub_rate"` // Set to DB in handle mode from audit-it.ru
}

const (
  UnitMeasure = 5000
  Currency = "RUB"
  EmptyDBErrorMessage = "DB is empty"
)

func (self *Main) GetHandler() map[string]interface{} {
  self.formExpenses()

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
)

func (self *Main) formExpenses() {
  dbExpenses := []DBExpense{}
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
      date := dbExpense.Date
      timestamp := int(date.Unix())
      year := date.Year()
      month := date.Month()
      day := date.Day()

      if month != prevDate.Month() {
        formatedMonths := self.addEmptyExpenses(expensesMonth, APIExpensesMonth, prevDate, true, weekNumber)

        // Behavior with empty month between fill months doesn't provided
        expensesYear = append(expensesYear, formatedMonths[0])
        APIExpensesYear = append(APIExpensesYear, formatedMonths[1])

        if year != prevDate.Year() {
          self.Expenses = append(self.Expenses, expensesYear)
          self.APIExpenses = append(self.APIExpenses, APIExpensesYear)

          expensesYear = [][]map[string]interface{}{}
          APIExpensesYear = [][]map[string]interface{}{}
        }

        firstMonthDay :=
          time.Unix(
            int64(timestamp - (day - 1) * DayTimestamp),
            0,
          )

        monthOffset = int(firstMonthDay.Weekday()) - 1 // Sunday is last weekday in EU

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
        weekNumber = (monthOffset + day) / DaysInWeek + 1
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
          "year_average_usd_rub_rate": dbExpense.YearAverageUSDRUBRate,
        }

      if dbExpense.YearAverageUSDRUBRate > 0 {
        if commentLength > 0 {
          APIExpense =
            map[string]interface{}{
              "id": dbExpense.Id,
              "value": dbExpense.Value,
              "comment": dbExpense.Comment,
              "year_average_usd_rub_rate": dbExpense.YearAverageUSDRUBRate,
            }
        } else {
          APIExpense =
            map[string]interface{}{
              "id": dbExpense.Id,
              "value": dbExpense.Value,
              "year_average_usd_rub_rate": dbExpense.YearAverageUSDRUBRate,
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

      if month == prevDate.Month() {
        gap := weekNumber - prevWeekNumber

        for itr := 1; itr < gap; itr++ {
          emptyDate :=
            time.Unix(
              int64(
                timestamp - (gap - itr) * WeekTimestamp,
              ),
              0,
            )

          id := bson.NewObjectId()

          expensesMonth =
            append(
              expensesMonth,
              map[string]interface{}{
                "id": id,
                "date": emptyDate,
              },
            )

          APIExpensesMonth =
            append(
              APIExpensesMonth,
              map[string]interface{}{
                "id": id,
              },
            )
        }
      } else {
        for itr := 1; itr < weekNumber; itr++ {
          emptyDate :=
            time.Unix(
              int64(
                timestamp - (weekNumber - itr) * WeekTimestamp,
              ),
              0,
            )

          id := bson.NewObjectId()

          if emptyDate.Month() != date.Month() {
            emptyDate =
              time.Unix(
                int64(
                  timestamp - (day - 1) * DayTimestamp,
                ),
                0,
              )
          }

          expensesMonth =
            append(
              expensesMonth,
              map[string]interface{}{
                "id": id,
                "date": emptyDate,
              },
            )

          APIExpensesMonth =
            append(
              APIExpensesMonth,
              map[string]interface{}{
                "id": id,
              },
            )
        }
      }

      expensesMonth = append(expensesMonth, expense)
      APIExpensesMonth = append(APIExpensesMonth, APIExpense)

      prevDate = date
      prevWeekNumber = weekNumber

      if key + 1 == dbExpensesLength {
        // Non full year
        formatedMonths := self.addEmptyExpenses(expensesMonth, APIExpensesMonth, prevDate, false, weekNumber)

        expensesYear = append(expensesYear, formatedMonths[0])
        APIExpensesYear = append(APIExpensesYear, formatedMonths[1])

        var yearIsAlreadyClosed bool

        now := time.Now()
        currentMonth := now.Month()
        currentYear := now.Year()

        // Fill empty months or years
        if month != currentMonth || year != currentYear {
          monthInt := int(month)
          gap := int(currentMonth) - monthInt

          if month != currentMonth {
            for itr := 0; itr < gap; itr++ {
              extraMonths := self.addExtraMonths(date, itr)

              expensesYear = append(expensesYear, extraMonths[0])
              APIExpensesYear = append(APIExpensesYear, extraMonths[1])
            }

            self.Expenses = append(self.Expenses, expensesYear)
            self.APIExpenses = append(self.APIExpenses, APIExpensesYear)
          }

          if year != currentYear {
            extraYears := int(currentYear) - int(year)

            if gap < 0 {
              gap = gap + MonthsInYear
            }

            for itr := 0; itr < extraYears; itr++ {
              emptyYear := [][]map[string]interface{}{}
              APIEmptyYear := [][]map[string]interface{}{}

              for monthItr := 0; monthItr < gap; monthItr++ {
                extraMonths := self.addExtraMonths(date, monthItr)

                emptyYear = append(emptyYear, extraMonths[0])
                APIEmptyYear = append(APIEmptyYear, extraMonths[1])
              }

              self.Expenses = append(self.Expenses, emptyYear)
              self.APIExpenses = append(self.APIExpenses, APIEmptyYear)
            }
          }

          yearIsAlreadyClosed = true
        }

        if yearIsAlreadyClosed != true {
          self.Expenses = append(self.Expenses, expensesYear)
          self.APIExpenses = append(self.APIExpenses, APIExpensesYear)
        }
      }
    }
  }
}

func (self *Main) addEmptyExpenses(
  month []map[string]interface{}, APIMonth []map[string]interface{},
  prevDate time.Time, redefineWeek bool, weekNumber int,
  ) [2][]map[string]interface{} {
    addition := WeeksInMonth - weekNumber

    if addition > 0 {
      for itr := 0; itr < addition; itr++ {
        id := bson.NewObjectId()
        date :=
          time.Unix(
            int64(
              int(prevDate.Unix()) + (itr + 1) * WeekTimestamp,
            ),
            0,
          )

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

func (self *Main) addExtraMonths(date time.Time, parentItr int) [2][]map[string]interface{} {
  month := []map[string]interface{}{}
  APIMonth := []map[string]interface{}{}

  for itr := 0; itr < WeeksInMonth; itr++ {
    var extraDate time.Time

    id := bson.NewObjectId()
    day := date.Day() - 1

    firstDayOfMonth :=
      time.Unix(
        int64(
          int(date.Unix()) - day * DayTimestamp,
        ),
        0,
      )

    firstDayOfMonth = firstDayOfMonth.AddDate(0, parentItr + 1, 0)

    if itr > 0 {
      firstDayOfMonthWeekday := int(firstDayOfMonth.Weekday())

      // Sunday is last weekday in EU
      if firstDayOfMonthWeekday == 0 {
        firstDayOfMonthWeekday = DaysInWeek
      }

      firstDayOfMonthTimestamp := int(firstDayOfMonth.Unix())
      daysInFirstWeek := DaysInWeek - firstDayOfMonthWeekday + 1

      extraDate =
        time.Unix(
          int64(firstDayOfMonthTimestamp + daysInFirstWeek * DayTimestamp + (itr - 1) * WeekTimestamp),
          0,
        )
    } else {
      extraDate = firstDayOfMonth
    }

    month =
      append(
        month,
        map[string]interface{}{
          "id": id,
          "date": extraDate,
        },
      )

    APIMonth =
      append(
        APIMonth,
        map[string]interface{}{
          "id": id,
        },
      )
  }

  result := [2][]map[string]interface{}{}

  result[0] = month
  result[1] = APIMonth

  return result
}

type SetReq struct {
  Id string `json: "id"`
  Value int `json: "value"`
  Comment string `json: "comment"`
}

func (self *Main) SetHandler(res *http.Request) map[string]interface{} {
  reqExpense := SetReq{}

  json.Unmarshal(
    []byte(
      self.ProcessReqBody(res),
    ),
    &reqExpense,
  )

  if reqExpense.Value != 0 && len(reqExpense.Id) > 0 {
    var value int

    dbExpense := DBExpense{}

    self.MongoCollection.Find(
      bson.M{
        "_id": bson.ObjectIdHex(reqExpense.Id),
      },
    ).One(&dbExpense)

    if len(dbExpense.Id) > 0 {
      value = dbExpense.Value + reqExpense.Value

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

      for _, year := range self.Expenses {
        for _, month := range year {
          for _, expense := range month {
            expenseId := expense["id"].(bson.ObjectId)

            if bson.ObjectId.Hex(expenseId) == reqExpense.Id {
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