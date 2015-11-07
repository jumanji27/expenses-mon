export default class Expense extends Backbone.View {
  constructor(model) {
    super({
      el: '.js_l_main',
      events: {
        'click .js_popup__add': 'set',
        'click .js_popup__remove': 'remove'
      }
    });

    this.model = model;
  }


  render(target, expense) {
    target.append(
      tmpl_components_shared_expense_main({
        id: expense.id || null,
        value: expense.value,
        amount: expense.amount
      })
    );
  }

  set() {
    this.model.addReq(
      $(this.el).find('.js_popup-start_active').attr('data-id')
    );
  }

  remove() {
    this.model.removeReq(
      $(this.el).find('.js_popup-start_active').attr('data-id')
    );
  }
}