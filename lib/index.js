/**
 * lib/index.js
 * Does the logic behind the hooks
 */

const debug = require('debug')
const util = require('util')
const path = require('path')

const helpers = require('./helpers')
const server = require('./server')

const listen = server.server
const log = debug('gits')

/**
 * setup
 * Sets up a new / updates an exisiting repository
 *
 * @param {Object} p - project configuration object
 * @param {String} b - branch to checkout
 */
const setup = (p, b) => {
  return helpers.gitto(p, b)
  .then(() => {
    if (p.script) {
      return helpers.exec(path.join('..', p.script))
    }
  })
}

/**
 * teardown
 * Removes a repository
 *
 * @param {Object} p - project configuration object
 * @param {String} b - branch to remove
 */
const teardown = (p, b) => {
  const pr = util.format(p.path, b)

  return helpers.rmr(pr)
}

listen.on('push', (res, project) => {
  const br = res.ref.split('/')[2]
  log(`Push occured for ${project.owner}/${project.repo}#${br}`)

  if (res.deleted) {
    teardown(project, br)
    .then(() => log(`Removed ${project.owner}/${project.repo}#${br}`))
    .catch((err) => {
      log(`Unable to remove ${project.owner}/${project.repo}#${br} due to error`)
      console.error(err)
    })

    return
  }

  setup(project, br)
  .then(() => log(`Updated ${project.owner}/${project.repo}#${br}`))
  .catch((err) => {
    log(`Unable to update ${project.owner}/${project.repo}#${br} due to error`)
    console.error(err)
  })
})

listen.on('create', (res, project) => {
  if (res.ref_type !== 'branch') return
  log(`New branch "${res.ref}" being created for ${project.owner}/${project.repo}`)

  setup(project, res.ref)
  .then(() => log(`Created ${project.owner}/${project.repo}#${res.ref}`))
  .catch((err) => {
    log(`Unable to create ${project.owner}/${project.repo}#${res.ref} due to error`)
    console.error(err)
  })
})
