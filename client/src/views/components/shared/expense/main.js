export default class Expense extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });
  }


  render(target, expense) {
    target.append(
      tmpl_components_shared_expense_main({
        value: expense.value,
        amount: expense.amount
      })
    );
  }
}

