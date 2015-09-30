export default class Expense extends Backbone.View {
  constructor(model) {
    super();

    this.model = model;

    this.listenTo(this.model, 'change', this.render);
  }

  render() {
    // Loop up data lvl here
    console.log(this.model);
  }
}