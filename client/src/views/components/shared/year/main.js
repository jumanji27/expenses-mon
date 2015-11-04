import Month from '../month/main';


export default class Year extends Backbone.View {
  constructor(expenses, renderTarget) {
    super();

    let self = this;

    expenses.map((year, key) => {
      self.render(renderTarget);

      new Month(
        year,
        key,
        $('.js_year').eq(key)
      );
    });
  }


  render(target) {
    target.append(tmpl_components_shared_year_main());
  }
}