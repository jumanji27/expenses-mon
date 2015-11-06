import Year from '../shared/year/main';
import Month from '../shared/month/main';
import Expense from '../shared/expense/main';


export default class Main extends Backbone.View {
  constructor(model, renderTarget) {
    super({
      el: '.js_wrapper'
    });

    this.model = model;

    this.listenTo(
      this.model,
      'change',
      () => {
        this.render(renderTarget);
      }
    );
  }


  render(target) {
    target.html(tmpl_components_main_main());

    let yearView = new Year(),
      monthView = new Month(),
      expenseView = new Expense(this.model),
      mainEl = $(this.el).find('.js_p-main');

    this.model.get('expenses').map((year, key) => {
      yearView.render(mainEl);

      year.map((month, monthKey) => {
        let yearEl = mainEl.children('.js_year').eq(key);

        monthView.render(
          yearEl,
          month[0].month
        );

        month.map((expense) => {
          expenseView.render(
            yearEl.children('.js_month').eq(monthKey),
            expense
          )
        });
      });
    });

    $(this.el).find('.js_popup-start').simplePopup();
  }
}