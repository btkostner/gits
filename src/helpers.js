/**
 * src/helpers.js
 * Some useful functions
 *
 * @exports {Function} mkdirp - creates a nested directory
 * @exports {Function} rmr - removes a directory recursively
 * @exports {Function} gitto - Gets git repository to the latest branch
 * @exports {Function} exec - runs a script in child process
 */

const child = require('child_process')
const git = require('nodegit')
const os = require('os')
const path = require('path')
const Promise = require('bluebird')

const fs = Promise.promisifyAll(require('fs'))

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
 * exec
 * Runs a command
 *
 * @param {String} c - the command to run
 * @param [String] p - the path run the command in
 */
const exec = (c, p = os.tmpdir()) => {
  return new Promise((resolve, reject) => {
    child.exec(c, {
      cwd: p
    }, (err, stdout, stderr) => {
      if (err) return reject(err)
      return resolve()
    })
  })
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
  return mkdirp(pa)
  .then(() => fs.statAsync(path.resolve(pa, '.git')))
  .catch({ code: 'ENOENT' }, () => git.Clone(`https://github.com/${pr.owner}/${pr.repo}.git`, pa, {
    checkoutBranch: br
  }))
  .then(() => git.Repository.open(pa))
  .then((repo) => {
    return repo.fetch('origin')
    .then(() => repo.mergeBranches(br, `origin/${br}`))
  })
}

module.exports = {
  mkdirp,
  rmr,
  gitto,
  exec
}
