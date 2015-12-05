export default class Month extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });

    let total = $('.js_year__total')

    console.log(total);
  }


  render(target, month) {
    target.append(
      tmpl_components_shared_month_main({
        month: month
      })
    );


  }
}