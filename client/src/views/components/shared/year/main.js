import Month from '../month/main';


export default class Year extends Backbone.View {
  constructor(expenses, renderTarget) {
    super({
      el: '.js_p-main'
    });

    let self = this;
    let month = new Month();

    expenses.map((year, key) => {
      self.render(renderTarget);

      month.renderAll(
        $(self.el).children('.js_year').eq(key),
        year,
        key
      );
    });
  }


  render(target) {
    target.append(tmpl_components_shared_year_main());
  }
}