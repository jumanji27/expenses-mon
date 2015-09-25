export default class Router extends Backbone.Router {
  routes: {
    '': 'main',
    'about': 'about'
  }

  initialize() {
    Backbone.history.start();
  }

  main() {
    console.log('hello main');
  }
}