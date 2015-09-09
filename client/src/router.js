import $ from 'jquery';
import Backbone from 'backbone';


export default Backbone.Router.extend({
  routes: {
    '': 'main',
  },

  initialize: () => {
    console.log('router!');
  },

  main: () => {
    console.log('main!');
  },

});


let router = new Router();

Backbone.history.start();