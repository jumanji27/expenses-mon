export default class Expense extends Backbone.View {
  constructor(model) {
    super({
      el: '.js_l_main'
    });

    this.model = model;
  }


  render(target, expense) {
    target.append(
      tmpl_components_shared_expense_main({
        id: expense.id,
        value: expense.value,
        amount: expense.amount,
        date: expense.date
      })
    );
  }

  updateHTML(params) {
    let expense = $(this.el).find('.js_popup-start_active'),
      value =
        parseInt(
          expense.attr('data-value')
        );

    if (params && params.decrement) {
      value--;
    } else {
      value++;
    }

    let rawAmount = value * this.model.get('unitMeasure'),
      amount = rawAmount.toString().replace(/000$/g, 'k'),
      amountEl = expense.children('.js_expense__amount');

    expense.attr('data-value', value);

    if (value) {
      amountEl.text(amount);
    } else {
      amountEl.text('');
    }
  }
}