/**
 * projects/static.js
 * A basic file for use in static sites
 */

const debug = require('debug')
const path = require('path')

const helpers = require('../src/helpers')
const log = debug('gits:php')

/**
 * push
 * Updates repository on new push
 *
 * @param {Object} res - GitHub hook object
 * @param {Object} project - the project configuration object
 * @param {String} br - parsed branch name
 */
module.exports['push'] = (res, project, br) => {
  const p = path.resolve(project.folder, br)

  log(`Updating ${project.owner}/${project.repo}#${br}`)
  return helpers.gitto(p, project, br)
  .then(() => log(`Updating ${project.owner}/${project.repo}#${br} complete`))
}

/**
 * create
 * Creates a new repository branch
 *
 * @param {Object} res - GitHub hook object
 * @param {Object} project - the project configuration object
 * @param {String} br - parsed branch name
 */
module.exports['create'] = (res, project, br) => {
  const p = path.resolve(project.folder, br)

  log(`Creating ${project.owner}/${project.repo}#${br}`)
  return helpers.mkdirp(p)
  .then(() => helpers.gitto(p, project, br))
  .then(() => log(`Creating ${project.owner}/${project.repo}#${br} complete`))
}

/**
 * delete
 * Removes a repository on branch deletion
 *
 * @param {Object} res - GitHub hook object
 * @param {Object} project - the project configuration object
 * @param {String} br - parsed branch name
 */
module.exports['delete'] = (res, project, br) => {
  const p = path.resolve(project.folder, br)

  log(`Deleting ${project.owner}/${project.repo}#${br}`)
  return helpers.rmr(p)
  .then(() => log(`Deleting ${project.owner}/${project.repo}#${br} complete`))
}
