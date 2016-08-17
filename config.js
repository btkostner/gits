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
    owner: 'elementary',
    repo: 'mvp',
    secret: null,
    folder: '/tmp/mvp',
    type: 'npm-build'
  }]
}
