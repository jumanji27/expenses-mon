import Year from '../../shared/year/main'
import Month from '../../shared/month/main'
import Expense from '../../shared/expense/main'


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
    target.html(tmpl_components_pages_home_main());

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
      comment = this.el.find('.js_popup__comment'),
      commentValue = comment.val();

    if (commentValue.length > 0) {
      params.forReq.comment = commentValue;
    }

    this.model.setReq(params);

    comment.val('');
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

    this.el.find('.js_popup__comment').val('');
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

    if (!this.statusTextAnimationInProcess) {
      if (typeof args.text === 'object') {
        const REPLACE_SYMBOL = '#';

        let expenseAmount =
          this.el.find('.js_popup-start_active').text();

        status.text(
          args.text.reduce((previousChunk, currentChunk) => {
            if (currentChunk === REPLACE_SYMBOL) {
              if (expenseAmount.length) {
                return `${previousChunk} (${expenseAmount})`;
              } else {
                return previousChunk;
              }
            }
            return previousChunk + currentChunk;
          })
        )
      } else {
        status.text(args.text);
      }

      status.addClass('js_popup__status-fade');

      const CORRECTED_TRANSITION_OPACITY = 700;

      this.statusTextAnimationInProcess = true;

      setTimeout(
        () => {
          status.text('');
          status.removeClass('js_popup__status-fade');
          this.statusTextAnimationInProcess = false;
        },
        CORRECTED_TRANSITION_OPACITY
      )
    }
  }
}