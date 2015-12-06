export default class Year extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });
  }


  render(params) {
    let valueSum = 0,
      averageUSDRUBrate,
      total;

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

    if (averageUSDRUBrate) {
      let totalRUB = valueSum * params.unitMeasure,
        totalAmount = totalRUB.toString().replace(/000$/g, 'k'),
        totalUSD = (totalRUB / averageUSDRUBrate).toString().replace(/000$/g, 'k');

      total = totalAmount + ' / ' + averageUSDRUBrate + ' = ' + totalUSD;
    }

    params.target.append(
      tmpl_components_shared_year_main({
        total: total
      })
    );
  }
}