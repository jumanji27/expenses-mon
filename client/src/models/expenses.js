export default class Expenses extends Backbone.Model {
  constructor() {
    super();

    this.req();
  }

  req() {
    $.ajax({
      type: 'POST',
      url: 'api/v1/get',
      success: (res) => {
        this.set({
          expenses: this.format(res.success.expenses)
        });
      }
    });
  }

  format(apiExpenses) {
    let WEEKS_IN_MONTH = 5,
      MONTHS = [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December"
      ];

    let expenses = [];

    for (let key in apiExpenses) {
      let year = [];

      for (let monthKey in apiExpenses[key]) {
        let month = [],
          prevWeek = 0

        for (let expenseKey in apiExpenses[key][monthKey]) {
          let weekGap = apiExpenses[key][monthKey][expenseKey].week - prevWeek;

          if (weekGap > 1) {
            for (let gapKey = 1; gapKey < weekGap; gapKey++) {
              month.push({
                value: 0
              });
            }
          }

          month.push({
            value: apiExpenses[key][monthKey][expenseKey].value
          });

          if (apiExpenses[key][monthKey].length === (parseInt(expenseKey) + 1) &&
            apiExpenses[key][monthKey][expenseKey].week !== WEEKS_IN_MONTH) {
              for (let lastMonthKey = 1; lastMonthKey <= WEEKS_IN_MONTH - apiExpenses[key][monthKey][expenseKey].week;
                lastMonthKey++) {
                  month.push({
                    value: 0
                  });
              }
          }
          month[0].month = MONTHS[monthKey];

          prevWeek = apiExpenses[key][monthKey][expenseKey].week;
        }

        year.push(month);

        let keyInt = parseInt(key);
        let monthKeyInt = parseInt(monthKey);

        if (keyInt + 1 === apiExpenses.length && monthKeyInt + 1 === apiExpenses[key].length) {
          let currentMonth = new Date().getMonth() - 1;
          let monthKeyInt = parseInt(monthKey);

          if (currentMonth > monthKeyInt) {
            for (let additionMonthsKey = 0; additionMonthsKey < currentMonth - monthKeyInt; additionMonthsKey++) {
              let emptyMonth = [];

              for (let additionMonthsExpenseKey = 1; additionMonthsExpenseKey < 5; additionMonthsExpenseKey++) {
                  emptyMonth.push({
                    value: 0,
                    week: additionMonthsExpenseKey
                  });
              }

              emptyMonth[0].month = MONTHS[monthKeyInt + additionMonthsKey + 1];

              year.push(emptyMonth);
            }

          }
        }
      }

      expenses.push(year);
    }

    return expenses.reverse();
  }
}