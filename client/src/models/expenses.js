export default class Expenses extends Backbone.Model {
  constructor() {
    super();

    this.API_HTTP_METHOD = 'POST';
    this.API_VERSION = '1';
    this.API_URL = '/api/v' + this.API_VERSION + '/';

    this.getReq();
  }


  getReq() {
    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'get',
      success: (res) => {
        this.set({
          expenses: this.format(res.success.expenses, res.success.unit_measure)
        });
      }
    });
  }

  setReq(params) {
    let self = this;

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'set',
      data: params.forReq,
      success: (res) => {
        self.updateViewStatus(params.view, res);
      }
    });
  }

  removeReq(params) {
    let self = this;

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'remove',
      data: {
        id: params.id
      },
      success: (res) => {
        self.updateViewStatus(params.view, res);
      }
    });
  }

  updateViewStatus(view, res) {
    let params = {
      success: res.success
    }

    if (res.success) {
      params.text = 'Success!';
    } else {
      params.text = res.error;
    }

    view.updateStatus(params);
  }

  format(dbExpenses, unitMeasure) {
    let WEEKS_IN_MONTH = 5,
      MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

    let expenses = [];

    dbExpenses.map((rawYear, key) => {
      let year = [];

      rawYear.map((rawMonth, monthKey) => {
        let month = [],
          prevWeek = 0

        rawMonth.map((expense, expenseKey) => {
          let weekGap = expense.week - prevWeek;

          if (weekGap > 1) {
            for (let gapKey = 1; gapKey < weekGap; gapKey++) {
              month.push({
                value: 0
              });
            }
          }

          let rawAmount = expense.value * unitMeasure,
            amount = rawAmount.toString().replace(/000$/g, 'k');

          month.push({
            id: expense.id,
            value: expense.value,
            amount: amount
          });

          if (rawMonth.length === expenseKey + 1 && expense.week !== WEEKS_IN_MONTH) {
            for (let lastMonthKey = 1; lastMonthKey <= WEEKS_IN_MONTH - expense.week; lastMonthKey++) {
              month.push({
                value: 0
              });
            }
          }
          month[0].month = MONTHS[monthKey];

          prevWeek = expense.week;
        });

        year.push(month);

        if (key + 1 === dbExpenses.length && monthKey + 1 === rawYear.length) {
          let currentMonth = new Date().getMonth();

          if (currentMonth > monthKey) {
            for (let additionMonthsKey = 0; additionMonthsKey < currentMonth - monthKey; additionMonthsKey++) {
              let emptyMonth = [];

              for (let additionMonthsExpenseKey = 0; additionMonthsExpenseKey < 5; additionMonthsExpenseKey++) {
                emptyMonth.push({
                  value: 0,
                  week: additionMonthsExpenseKey
                });
              }

              emptyMonth[0].month = MONTHS[monthKey + additionMonthsKey + 1];

              year.push(emptyMonth);
            }

          }
        }
      });

      expenses.push(year);
    });

    return expenses.reverse();
  }
}