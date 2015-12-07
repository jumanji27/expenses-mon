export default class Expenses extends Backbone.Model {
  constructor() {
    super();

    this.API_HTTP_METHOD = 'POST';
    this.API_VERSION = '1';
    this.API_URL = '/api/v' + this.API_VERSION + '/';

    this.getReq();
  }


  getReq() {
    let that = this;

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'get',
      success: (res) => {
        that.set({
          unitMeasure: res.success.unit_measure
        });
        that.set({
          expenses: that.format(that, res.success.expenses)
        });
      }
    });
  }

  setReq(args) {
    let that = this;

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'set',
      data: JSON.stringify(args.forReq),
      success: (res) => {
        that.sendStatusToView(args.page, res);

        // Update views instead of model — bad design for scaling!
        args.expenseView.updateHTML(args.forReq.value);
        args.yearView.updateTotal(args.yearId, args.forReq.value);
      }
    });
  }

  sendStatusToView(view, res) {
    let params = {
      success: res.success
    }

    if (res.success) {
      params.text = 'Success!';
    } else {
      params.text = res.error;
    }

    view.popupUpdateStatus(params);
  }

  format(that, dbExpenses) {
    let MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
      expenses = [];

    dbExpenses.map((rawYear, key) => {
      let year = [];

      rawYear.map((rawMonth, monthKey) => {
        let month = [];

        rawMonth.map((apiExpense, expenseKey) => {
          let expense = {
            id: apiExpense.id,
            value: apiExpense.value,
            yearAverageUSDRUBRate: apiExpense.year_average_usd_rub_rate
          }

          if (apiExpense.value) {
            let rawAmount = apiExpense.value * that.get('unitMeasure'),
              amount = rawAmount.toString().replace(/000$/g, 'k');

            expense.amount = amount;
          }

          month.push(expense);

          month[0].month = MONTHS[monthKey];
        });

        year.push(month);
      });

      expenses.push(year);
    });

    return expenses;
  }
}