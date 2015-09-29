export default class Expense extends Backbone.Model {
  constructor() {
    super();
  }

  req() {
    $.ajax(
      type: 'POST',
      url: 'api/v1/get',
      success: function(res) {
        console.log(res);
      }
    );
  }
}