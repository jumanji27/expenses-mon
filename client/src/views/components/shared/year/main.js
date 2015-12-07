export default class Year extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });
  }


  render(params) {
    params.target.append(
      tmpl_components_shared_year_main({
        id: params.id,
        total: this.getTotal(params)
      })
    );
  }

  updateTotal(id, value) {
    console.log($('.js_year[data-id=' + id + ']').children('.js_year__total'));

    $('.js_year[data-id=' + id + ']').children('.js_year__total').text(
      this.getTotal({
        updateValue: value
      })
    );
  }

  getTotal(params) {
    let valueSum = 0,
      averageUSDRUBrate;

    if (params.expenses) {
      this.reserveParams = params;
    } else {
      params = this.reserveParams;
    }

    params.expenses.map((month) => {
      month.map((expense) => {
        if (expense.value) {
          valueSum += expense.value;
        }

        if (expense.yearAverageUSDRUBRate) {
          averageUSDRUBrate = expense.yearAverageUSDRUBRate;
        }
      });
    });

    let updateValue = params.updateValue || 0,
      totalRUB = (valueSum + updateValue) * params.unitMeasure,
      total = totalRUB.toString().replace(/000$/g, 'k');

    if (averageUSDRUBrate) {
      let totalUSD = (totalRUB / averageUSDRUBrate).toString().replace(/000$/g, 'k');

      return total + ' / ' + averageUSDRUBrate + ' = $' + Math.round(totalUSD / 100) / 10 + 'k';
    } else {
      return total;
    }
  }
}