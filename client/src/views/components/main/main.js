import Expense from '../../../views/components/shared/expense/main';


export default class Main extends Backbone.View {
  constructor(model) {
    super();

    this.model = model;

    this.listenTo(this.model, 'change', this.render);
  }

  render() {
    let expenses = this.model.get('expenses');
    let html = '';
    let expense = new Expense()

    for (let year of expenses) {
      for (let month of year) {
        for (let expenseFromModel of month) {
          html += expense.returnHTML(expenseFromModel.value);
        }
      }
    }

    $('.js_wrapper').html(html);
  }
}