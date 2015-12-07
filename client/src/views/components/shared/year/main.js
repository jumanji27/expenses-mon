export default class Year extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });

    this.reserveParams = [];
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
      averageUSDRUBrate,
      updateValue = args.updateValue || 0;

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

    let totalRUB = (valueSum + updateValue) * args.unitMeasure,
      total = totalRUB.toString().replace(/000$/g, 'k');

    if (averageUSDRUBrate) {
      let totalUSD = (totalRUB / averageUSDRUBrate).toString().replace(/000$/g, 'k');

      return total + ' / ' + averageUSDRUBrate + ' = $' + Math.round(totalUSD / 100) / 10 + 'k';
    } else {
      return total;
    }
  }
}