import Year from '../shared/year/main';
import Month from '../shared/month/main';
import Expense from '../shared/expense/main';


export default class Main extends Backbone.View {
  constructor(model, renderTarget) {
    super({
      el: '.js_wrapper',
      events: {
        'click .js_popup__add': 'popupAdd',
        'click .js_popup__remove': 'popupRemove'
      }
    });

    this.model = model;

    this.listenTo(
      this.model,
      'change:expenses',
      () => {
        this.render(renderTarget);
      }
    );
  }


  render(target) {
    target.html(tmpl_components_main_main());

    let yearView = new Year(),
      monthView = new Month(),
      mainEl = $(this.el).find('.js_p-main');

    this.expenseView = new Expense(this.model);

    this.model.get('expenses').map((year, key) => {
      yearView.render(mainEl);

      year.map((month, monthKey) => {
        let yearEl = mainEl.children('.js_year').eq(key);

        monthView.render(
          yearEl,
          month[0].month
        );

        month.map((expense) => {
          this.expenseView.render(
            yearEl.children('.js_month').eq(monthKey),
            expense
          )
        });
      });
    });

    $(this.el).find('.js_popup-start').simplePopup();
  }

  popupAdd() {
    let expense = $(this.el).find('.js_popup-start_active');

    let params = {
      page: this,
      view: this.expenseView,
      forReq: {
        value: 1,
        id: expense.attr('data-id')
      }
    },
      comment = $(this.el).find('.js_popup__comment').val();

    if (comment.length > 0) {
      params.forReq.comment = comment;
    }

    this.model.setReq(params);
  }

  popupRemove() {
    let expense = $(this.el).find('.js_popup-start_active'),
      value =
        parseInt(
          expense.attr('data-value')
        );

    if (value) {
      this.model.setReq({
        page: this,
        view: this.expenseView,
        forReq: {
          value: -1,
          id: expense.attr('data-id')
        }
      });
    }
  }

  popupUpdateStatus(params) {
    let status = $(this.el).find('.js_popup__status'),
      statusHasErrorClass = status.hasClass('js_popup__status-error');

    if (params.success) {
      if (statusHasErrorClass) {
        status.removeClass('js_popup__status-error');
      }
    } else if (!statusHasErrorClass) {
      status.addClass('js_popup__status-error');
    }

    status.text(params.text);
  }
}