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
    this.el = $(this.el);

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

    let monthView = new Month(),
      mainEl = this.el.find('.js_p-main');

    this.yearView = new Year();
    this.expenseView = new Expense(this.model);

    this.model.get('expenses').map((year, key) => {
      this.yearView.render({
        target: mainEl,
        id: key + 1,
        expenses: year,
        unitMeasure: this.model.get('unitMeasure')
      });

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

    this.el.find('.js_popup-start').simplePopup();
  }

  popupAdd() {
    let expense = this.el.find('.js_popup-start_active'),
      params = {
        page: this,
        yearView: this.yearView,
        yearId: expense.closest('.js_year').attr('data-id'),
        expenseView: this.expenseView,
        forReq: {
          value: 1,
          id: expense.attr('data-id')
        }
      },
      comment = this.el.find('.js_popup__comment').val();

    if (comment.length > 0) {
      params.forReq.comment = comment;
    }

    this.model.setReq(params);
  }

  popupRemove() {
    let expense = this.el.find('.js_popup-start_active'),
      value =
        parseInt(
          expense.attr('data-value')
        );

    if (value) {
      this.model.setReq({
        page: this,
        yearView: this.yearView,
        yearId: expense.closest('.js_year').attr('data-id'),
        expenseView: this.expenseView,
        forReq: {
          value: -1,
          id: expense.attr('data-id')
        }
      });
    }
  }

  popupUpdateStatus(args) {
    let status = this.el.find('.js_popup__status'),
      statusHasErrorClass = status.hasClass('js_popup__status-error');

    if (args.success) {
      if (statusHasErrorClass) {
        status.removeClass('js_popup__status-error');
      }
    } else if (!statusHasErrorClass) {
      status.addClass('js_popup__status-error');
    }

    status.text(args.text);
  }
}