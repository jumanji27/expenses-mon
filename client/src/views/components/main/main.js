export default class Main extends Backbone.View {
  constructor(model) {
    super();

    this.model = model;

    this.listenTo(this.model, 'change', this.render);
  }

  render() {
    $('.js_wrapper').html(
      tmpl_components_main_main({
        expenses: this.model.get('expenses')
      })
    );
  }
}