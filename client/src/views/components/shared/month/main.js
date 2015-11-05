export default class Month extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
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