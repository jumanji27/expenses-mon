export default class Expense extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });
  }


  renderAll(target, month) {
    let self = this;

    month.map((expense, key) => {
      self.render(target, expense);
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

