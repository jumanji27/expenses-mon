import ExpenseModel from './models/expense'
import ExpenseView from './views/components/shared/expense/main'


export default class Router extends Backbone.Router {
  constructor() {
    super();

    this.routes = {
      '': 'main'
    }

    this.expense_mon = {
      collections: {},
      models: {},
      views: {}
    }

    this._bindRoutes();
    Backbone.history.start();
  }

  main() {
    this.expense_mon.models.expense = new ExpenseModel();
    this.expense_mon.views.expense = new ExpenseView(this.expense_mon.models.expense);
  }
}