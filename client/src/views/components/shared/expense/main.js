export default class Expense extends Backbone.View {
  constructor(model) {
    super({
      el: '.js_p_main'
    });

    this.model = model;
  }


  render(target, expense) {
    target.append(
      tmpl_components_shared_expense_main({
        id: expense.id || null,
        value: expense.value,
        amount: expense.amount || null
      })
    );
  }

  updateHTML(params) {
    let value =
      parseInt(
        params.expense.attr('data-value')
      );

    if (!params.decrement) {
      value++;
    } else {
      value--;
    }

    let rawAmount = value * this.model.get('unitMeasure'),
      amount = rawAmount.toString().replace(/000$/g, 'k');

    params.expense.attr('data-value', value);
    params.expense.children('.js_expense__amount').text(amount);
  }
}