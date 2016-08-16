/**
* config.js
* Holds all configuration variables for gits
*
* @exports {Object} - configuration
*/

module.exports = {
  port: 3000,
  host: '0.0.0.0',

  projects: [{
    owner: 'btkostner',
    repo: 'mvp',
    secret: null,
    folder: '/tmp/mvp',
    type: 'static'
  }]
}
