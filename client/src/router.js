import ExpensesModel from './models/expenses';
import MainPageView from './views/components/pages/home/main';


export default class Router extends Backbone.Router {
  constructor() {
    super();

    this.routes = {
      '': 'main'
    }

    this.expense_mon = {
      collections: {},
      models: {},
      views: {
        main: {}
      }
    }

    this._bindRoutes();
    Backbone.history.start();
  }

  main() {
    this.expense_mon.models.expenses = new ExpensesModel();
    this.expense_mon.views.main.main =
      new MainPageView(
        this.expense_mon.models.expenses,
        $('.js_wrapper')
      );
  }
}