/**
 * projects/npm-build.js
 * Runs `npm build` before serving. Inherits from static
 */

const debug = require('debug')
const path = require('path')

const helpers = require('../src/helpers')
const log = debug('gits:npm-build')

const stat = require('./static')

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

  return stat.push(res, project, br)
  .then(() => log(`Running 'npm install' on ${project.owner}/${project.repo}#${br}`))
  .then(() => helpers.exec('npm install', p))
  .then(() => log(`Running 'npm install' on ${project.owner}/${project.repo}#${br} complete`))
  .then(() => log(`Running 'npm run build' on ${project.owner}/${project.repo}#${br}`))
  .then(() => helpers.exec('npm run build', p))
  .then(() => log(`Running 'npm run build' on ${project.owner}/${project.repo}#${br} complete`))
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

  return stat.create(res, project, br)
  .then(() => log(`Running 'npm install' on ${project.owner}/${project.repo}#${br}`))
  .then(() => helpers.exec('npm install', p))
  .then(() => log(`Running 'npm install' on ${project.owner}/${project.repo}#${br} complete`))
  .then(() => log(`Running 'npm run build' on ${project.owner}/${project.repo}#${br}`))
  .then(() => helpers.exec('npm run build', p))
  .then(() => log(`Running 'npm run build' on ${project.owner}/${project.repo}#${br} complete`))
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
  return stat.create(res, project, br)
}
