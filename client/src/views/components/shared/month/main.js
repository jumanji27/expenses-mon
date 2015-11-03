export default class Month extends Backbone.View {
  constructor(expenses, renderTarget) {
    super();

    let self = this;

    expenses.map((expense, key) => {
      self.render(renderTarget);
    });
  }

  render(target) {
    console.log(target)

    // target.append(tmpl_components_shared_month_main());
  }
}