export default class Expense extends Backbone.View {
  constructor(month, renderTarget) {
    super();

    let self = this;

    month.map((expense, key) => {
      self.render(
        renderTarget,
        expense
      );
    });

    this.events = {
      'click .js_expense': 'openPopup'
    };
  }


  render(target, expense) {
    target.append(
      tmpl_components_shared_expense_main({
        value: expense.value,
        amount: expense.amount
      })
    );
  }

  openPopup() {
    console.log('xxx');
  }
}