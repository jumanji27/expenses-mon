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
    let WEEKS_IN_MONTH = 5;

    let expenses = [];

    for (let apiYear of apiExpenses) {
      let year = [];

      for (let apiMonth of apiYear) {
        let month = [],
          prevWeek = 0

        for (let key in apiMonth) {
          let weekGap = apiMonth[key].week - prevWeek;

          if (weekGap > 1) {
            for (let gapKey = 1; gapKey < weekGap; gapKey++) {
              month.push({
                value: 0
              });
            }
          }

          month.push({
            value: apiMonth[key].value
          });

          if (apiMonth.length === (parseInt(key) + 1) && apiMonth[key].week !== WEEKS_IN_MONTH) {
            for (let lastMonthKey = 1; lastMonthKey <= WEEKS_IN_MONTH - apiMonth[key].week; lastMonthKey++) {
              month.push({
                value: 0
              });
            }
          }

          prevWeek = apiMonth[key].week;
        }

        year.push(month)
      }

      expenses.push(year);
    }

    return expenses;
  }
}