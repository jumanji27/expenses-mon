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
          expenses: this.addEmptyWeeks(res.success.expenses)
        });
      }
    });
  }

  addEmptyWeeks(apiExpenses) {
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

        year.push(month)
      }

      expenses.push(year);
    }

    return expenses.reverse();
  }
}