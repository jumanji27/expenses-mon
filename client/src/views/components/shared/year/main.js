export default class Year extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });

    this.reserveParams = [];
    this.updateValue = 0;
  }


  render(args) {
    args.target.append(
      tmpl_components_shared_year_main({
        id: args.id,
        total: this.getTotal(args)
      })
    );
  }

  updateTotal(id, value) {
    $('.js_year[data-id="' + id + '"]').children('.js_year__total').text(
      this.getTotal({
        yearId: id,
        updateValue: value
      })
    );
  }

  getTotal(args) {
    let valueSum = 0,
      averageUSDRUBrate;

    if (args.updateValue || args.updateValue === 0) {
      this.updateValue += args.updateValue;
    }

    if (args.expenses) {
      this.reserveParams.push(args);
    } else {
      args = this.reserveParams[args.yearId - 1];
    }

    args.expenses.map((month) => {
      month.map((expense) => {
        if (expense.value) {
          valueSum += expense.value;
        }

        if (expense.yearAverageUSDRUBRate) {
          averageUSDRUBrate = expense.yearAverageUSDRUBRate;
        }
      });
    });

    let year = new Date().getFullYear() - args.id + 1,
      totalRUB = (valueSum + this.updateValue) * args.unitMeasure,
      total = totalRUB.toString().replace(/000$/g, 'k'),
      totalWithCurrency;

    if (averageUSDRUBrate) {
      let totalUSD = (totalRUB / averageUSDRUBrate).toString().replace(/000$/g, 'k');

      totalWithCurrency = total + ' / ' + averageUSDRUBrate + ' = $' + Math.round(totalUSD / 100) / 10 + 'k';
    } else {
      totalWithCurrency = total;
    }

    return year + ': ' + totalWithCurrency;
  }
}