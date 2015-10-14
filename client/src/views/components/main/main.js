import Month from '../../../views/components/shared/month/main';


export default class Main extends Backbone.View {
  constructor(model) {
    super();

    this.model = model;

    this.listenTo(this.model, 'change', this.render);
  }

  render() {
    let html = '',
      expenses = this.model.get('expenses'),
      month = new Month();

    for (let year of expenses) {
      for (let monthFromModel of year) {
        html += month.returnHTML(monthFromModel);
      }
    }

    $('.js_wrapper').html(html);
  }
}