class Router extends Backbone.Router {
  constructor () {
    super();
  }

  routes: {
    '': 'main'
  }

  main () {
    console.log('Router#home was called!');
  }
}

// export default Router;