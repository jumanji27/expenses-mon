export default class Expense extends Backbone.View {
  constructor() {
    super();
  }

  returnHTML(value) {
    return tmpl_components_shared_expense_main({
      value: value
    });
  }
}