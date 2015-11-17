export default class Expenses extends Backbone.Model {
  constructor() {
    super();

    this.API_HTTP_METHOD = 'POST';
    this.API_VERSION = '1';
    this.API_URL = '/api/v' + this.API_VERSION + '/';

    this.getReq();
  }


  getReq() {
    let self = this;

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'get',
      success: (res) => {
        self.set({
          unitMeasure: res.success.unit_measure
        });
        self.set({
          expenses: self.format(self, res.success.expenses)
        });
      }
    });
  }

  setReq(params) {
    let self = this;

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'set',
      data: JSON.stringify(params.forReq),
      success: (res) => {
        self.sendStatusToView(params.page, res);

        let paramsToView = {};

        if (params.forReq.action === 'remove') {
          paramsToView.decrement = true;
        }

        params.view.updateHTML(paramsToView);
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

  format(self, dbExpenses) {
    let MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
      expenses = [];

    dbExpenses.map((rawYear, key) => {
      let year = [];

      rawYear.map((rawMonth, monthKey) => {
        let month = [];

        rawMonth.map((apiExpense, expenseKey) => {
          let expense = {
            id: apiExpense.id,
            value: apiExpense.value,
            date: apiExpense.date
          }

          if (apiExpense.value) {
            let rawAmount = apiExpense.value * self.get('unitMeasure'),
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