'use strict';

var Promise = require('es6-promise').Promise;

var Prompt = require('../../cli/prompt');
var Command = require('../../cli/command');

var login = require('../../login');

module.exports = new Command(
  'login',
  'log in to your Arigato account',
  function(ctx) {
    return new Promise(function(resolve, reject) {

      // Create prompt from login questions
      var prompt = new Prompt(login.questions);

      // Begin asking questions
      prompt.start().then(function(userInput) {

        // Attempt login from user input
        return login.attempt(ctx.daemon, userInput);

      // Success, session created
      }).then(function() {
        // TODO: Proper output module for errors and banner messages
        console.log('');
        console.log('You are now authenticated');
        console.log('');
        resolve();

      // Account creation failed
      }).catch(function(err) {
        err.type = err.type || 'unknown';
        console.error('');
        console.error('Login failed, please try again');
        console.error('');
        reject(err);
      });
    });
  }
);