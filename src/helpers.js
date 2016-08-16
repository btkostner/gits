/**
 * src/helpers.js
 * Some useful functions
 *
 * @exports {Function} mkdirp - creates a nested directory
 * @exports {Function} rmr - removes a directory recursively
 * @exports {Function} gitto - Gets git repository to the latest branch
 * @exports {Function} exec - runs a script in child process
 */

const debug = require('debug')
const git = require('nodegit')
const path = require('path')
const Promise = require('bluebird')
const util = require('util')

const fs = Promise.promisifyAll(require('fs'))
const log = debug('gits:helpers')

/**
 * mkdirp
 * Creates a directory
 *
 * @param {String} p - path to folder to create
 * @returns {Void}
 */
const mkdirp = (p) => {
  const chunks = p.split(path.sep)

  return Promise.map(chunks, (chunk, i) => {
    return chunks.slice(0, i + 1).join(path.sep)
  })
  .each((chunk) => {
    if (chunk == null || chunk === '') return

    return fs.mkdirAsync(chunk)
    .catch({ code: 'EEXIST' }, () => true)
  })
}

/**
 * rmr
 * Removes a directory recursively
 *
 * @param {String} p - path to remove
 * @returns {Void}
 */
const rmr = (p) => {
  return fs.statAsync(p)
  .then((stat) => {
    if (stat.isFile()) {
      return fs.unlinkAsync(p)
    } else {
      return fs.readdirAsync(p)
      .each((i) => rmr(path.join(p, i)))
      .then(() => fs.rmdirAsync(p))
    }
  })
  .catch({ code: 'ENOENT' }, () => true)
}

/**
 * gitto
 * Gets git repository to the latest branch
 * FIXME: this leaves the repo in a detached head state
 *
 * @param {String} pa - full path of the new branch
 * @param {Object} pr - project configuration object
 * @param {String} br - name of branch to get to
 */
const gitto = (pa, pr, br) => {
  let repo = null

  return git.Repository.open(pa)
  .catch(() => {
    log(`Creating new git repository "${pr.repo}#${br}"`)
    return git.Clone(`https://github.com/${pr.owner}/${pr.repo}.git`, pa)
  })
  .then((r) => { repo = r })
  .then(() => repo.fetchAll())
  .then(() => repo.getReference(`origin/${br}`))
  .then((ref) => repo.checkoutRef(ref, {
    checkoutStrategy: git.Checkout.STRATEGY.USE_THEIRS
  }))
}

module.exports = {
  mkdirp,
  rmr,
  gitto
}
