import Year from '../shared/year/main';


export default class Main extends Backbone.View {
  constructor(model, renderTarget) {
    super();

    this.model = model;

    this.listenTo(
      this.model,
      'change',
      () => {
        this.render(renderTarget);
      }
    );
  }

  render(target) {
    target.html(tmpl_components_main_main());

    new Year(
      this.model.get('expenses'),
      $('.js_p-main')
    );
  }
}