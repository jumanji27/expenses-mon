import Expense from '../expense/main';


export default class Month extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });

    this.expense = new Expense();
  }


  renderAll(target, year, yearKey) {
    let self = this;

    year.map((month, key) => {
      self.render(target, month[0].month);

      this.expense.renderAll(
        $(self.el).children('.js_year').eq(yearKey).children('.js_month').eq(key),
        month
      );
    });
  }

  render(target, month) {
    target.append(
      tmpl_components_shared_month_main({
        month: month
      })
    );
  }
}