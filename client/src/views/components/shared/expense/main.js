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

  updateHTML(value) {
    let expense = $(this.el).find('.js_popup-start_active'),
      HTMLValue = expense.attr('data-value');

    if (HTMLValue) {
      HTMLValue = parseInt(HTMLValue);

      if (value > 0) {
        HTMLValue++;
      } else {
        HTMLValue--;
      }
    } else {
      HTMLValue = 1;
    }

    let rawAmount = HTMLValue * this.model.get('unitMeasure'),
      amount = rawAmount.toString().replace(/000$/g, 'k'),
      amountEl = expense.children('.js_expense__amount');

    expense.attr('data-value', HTMLValue);

    if (HTMLValue) {
      amountEl.text(amount);
    } else {
      amountEl.text('');
    }
  }
}