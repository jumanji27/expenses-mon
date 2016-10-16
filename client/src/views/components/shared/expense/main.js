export default class Expense extends Backbone.View {
  constructor(model) {
    super({
      el: '.js_l_main'
    });

    this.model = model;
    this.el = $(this.el);
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
    let expense = this.el.find('.js_popup-start_active'),
      DOMValue = expense.attr('data-value');

    if (DOMValue) {
      DOMValue = parseInt(DOMValue);

      if (value > 0) {
        DOMValue++;
      } else {
        DOMValue--;
      }
    } else {
      DOMValue = 1;
    }

    let amount = this.model.formatAmount(DOMValue),
      amountEl = expense.children('.js_expense__amount');

    expense.attr('data-value', DOMValue);

    if (DOMValue) {
      amountEl.text(amount);
    } else {
      amountEl.text('');
    }
  }
}