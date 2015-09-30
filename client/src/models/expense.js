export default class Expense extends Backbone.Model {
  constructor() {
    super();

    this.req();
  }

  req() {
    $.ajax({
      type: 'POST',
      url: 'api/v1/get',
      success: (res) => {
        this.set(res.success);
      }
    });
  }
}