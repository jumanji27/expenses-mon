export default class Expense extends Backbone.View {
  constructor(model) {
    super();

    this.model = model;

    this.listenTo(this.model, 'change', this.render);
  }

  render() {
    let expenses = this.model.get('expenses');
    let html = '';

    for (let year of expenses) {
      for (let month of year) {
        for (let expense of month) {
          console.log(expense);
        }
      }
    }


  }
}