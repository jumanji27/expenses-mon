import Year from '../shared/year/main';


export default class Main extends Backbone.View {
  constructor(model, renderTarget) {
    super({
      el: '.js_wrapper'
    });

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
      $(this.el).find('.js_p-main')
    );

    $(this.el).find('.js_popup-start').simplePopup();
  }
}