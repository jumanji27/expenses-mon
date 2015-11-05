export default class Year extends Backbone.View {
  constructor() {
    super({
      el: '.js_p-main'
    });
  }


  render(target) {
    target.append(tmpl_components_shared_year_main());
  }
}