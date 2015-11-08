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
    let expense = $(this.el).find('.js_popup-start_active'),
      value =
        parseInt(
          expense.attr('data-value')
        );

    expense.attr('data-value', value + 1);

    let params = {
      view: this,
      forReq: {
        id: expense.attr('data-id')
      }
    },
      comment = $(this.el).find('.js_popup__comment').val();

    if (comment.length > 0) {
      params.forReq.comment = comment;
    }

    this.model.setReq(params);
  }

  remove() {
    let expense = $(this.el).find('.js_popup-start_active'),
      value =
        parseInt(
          expense.attr('data-value')
        );

    expense.attr('data-value', value - 1);

    this.model.removeReq({
      view: this,
      id: expense.attr('data-id')
    });
  }

  updateStatus(params) {
    let status = $(this.el).find('.js_popup__status');

    if (!params.success) {
      status.addClass('js_popup__status-error')
    }

    status.text(params.text);
  }
}