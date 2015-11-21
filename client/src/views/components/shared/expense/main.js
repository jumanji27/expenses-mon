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
        amount: expense.amount
      })
    );
  }

  updateHTML(params) {
    let expense = $(this.el).find('.js_popup-start_active'),
      value = expense.attr('data-value');

    if (value) {
      value = parseInt(value);

      if (params && params.decrement) {
        value--;
      } else {
        value++;
      }
    } else {
      value = 1;
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