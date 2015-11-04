import Expense from '../expense/main';


export default class Month extends Backbone.View {
  constructor(year, yearKey, renderTarget) {
    super();

    let self = this;

    year.map((month, key) => {
      self.render(renderTarget, month[0].month);

      new Expense(
        month,
        $('.js_year').eq(yearKey).children('.js_month').eq(key)
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