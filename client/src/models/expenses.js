export default class Expenses extends Backbone.Model {
  constructor() {
    super();

    this.API_HTTP_METHOD = 'POST';
    this.API_VERSION = '1';
    this.API_URL = '/api/v' + this.API_VERSION + '/';

    this.getReq();
  }


  getReq() {
    const NAME = this.API_URL + 'get';

    $.ajax({
      type: this.API_HTTP_METHOD,
      url: NAME,
      success: (res) => {
        this.set({
          unitMeasure: res.success.unit_measure
        });
        this.set({
          expenses: this.format(res.success.expenses)
        });

        localStorage.setItem(
          NAME,
          JSON.stringify(res.success)
        );
      },
      error: () => {
        let cacheRes = localStorage.getItem(NAME),
          cachedResLength = cacheRes.length;

        if (cachedResLength) {
          this.set({
            unitMeasure: JSON.parse(cacheRes).unit_measure
          });
          this.set({
            expenses:
              this.format(
                JSON.parse(cacheRes).expenses
              )
          });
        }

        console.log('[LS] Working from cache')
      }
    });
  }

  setReq(args) {
    $.ajax({
      type: this.API_HTTP_METHOD,
      url: this.API_URL + 'set',
      data: JSON.stringify(args.forReq),
      success: (res) => {
        let params = {
          success: res.success
        }

        if (res.success) {
          // # — special symbol for replacements
          params.text = ['Success!', '#'];

          // Update views instead of model — bad design for scaling!
          args.expenseView.updateHTML(args.forReq.value);
          args.yearView.updateTotal(args.yearId, args.forReq.value);
        } else {
          params.text = res.error;
        }

        args.page.popupUpdateStatus(params);
      },
      error: () => {
        args.page.popupUpdateStatus({
          text: "You're offline, sorry"
        });

        console.log('[LS] Working from cache')
      }
    });
  }

  format(dbExpenses) {
    const MONTHS = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

    let expenses = [];

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
            let amount = this.formatAmount(apiExpense.value);

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

  formatAmount(value) {
    let rawAmount = value * this.get('unitMeasure');

    return rawAmount.toString().replace(/000$/g, 'k');
  }
}