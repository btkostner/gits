/**
* lib.server.js
* Starts the github hook server
*
* @exports {Object} http - http server already listening
* @exports {Object} server - server eventemitter
*/

const crypto = require('crypto')
const debug = require('debug')
const Event = require('events')
const http = require('http')

const config = require('../config.js')
const log = debug('gits:server')
const server = new Event()

/**
* reply
* Sends a reply back to GitHub
*
* @param {Object} res - http server response object
* @param {Number} status - http response code
* @param {String} msg - message to sendback to GitHub
* @returns {Void}
*/
const reply = (res, status, msg) => {
  const message = JSON.stringify({
    message: msg || http.STATUS_CODES[status] || null,
    result: (status <= 400) ? 'ok' : 'error'
  })

  const headers = {
    'Content-Type': 'application/json',
    'Content-Length': message.length
  }

  res.writeHead(status, headers)
  res.end(message)
}

/**
* handler
* Handles all http connections
*
* @param {Object} req - the data sent from browser
* @param {Object} res - the response we send back to the browser
* @returns {Void}
*/
const handler = (req, res) => {
  if (req.method !== 'POST') {
    log('Request is not a POST request')
    reply(res, 405)
    return
  }

  if (req.headers['x-github-event'] == null) {
    log('Request does not have a github event header')
    reply(res, 400)
    return
  }

  const buffer = []
  let length = 0

  req.on('data', (chunk) => {
    log('Request sent new chunk')

    buffer.push(chunk)
    length += chunk.length
  })

  req.on('end', (chunk) => {
    log('Request ended')

    if (chunk) {
      buffer.push(chunk)
      length += chunk.length
    }

    let data = null
    try {
      data = Buffer.concat(buffer, length).toString()
    } catch (err) {
      log('Error while trying to concat buffer')
      console.error(err)

      reply(res, 500)
      return
    }

    let pkg = null
    try {
      pkg = JSON.parse(data)
    } catch (e) {
      log('Unable to parse data')
      reply(res, 403, 'Unable to parse data')
      return
    }

    if (pkg == null || pkg.repository == null || pkg.repository.full_name == null) {
      reply(res, 403, 'Unable to identify repository name')
      return
    }

    const project = config.projects.find((project) => {
      return `${project.owner}/${project.repo}` === pkg.repository.full_name
    })
    if (project == null) {
      log(`Request was for ${pkg.repository.full_name}, and is not in configuration`)
      reply(res, 404, 'Repository not found')
      return
    }

    let signature = req.headers['x-hub-signature']
    if (signature == null) {
      log('Request did not send a signature. Assuming insecure')
      reply(res, 403, 'Sent without a signature')
      return
    }
    signature = signature.replace(/^sha1=/, '')

    const digest = crypto.createHmac('sha1', project.secret).update(data).digest('hex')
    if (signature !== digest) {
      log('Signature and digest do not match!')
      reply(res, 403)
      return
    }

    const type = req.headers['x-github-event']
    if (type == null) {
      log('Request did not send an event type')
      reply(res, 403, 'Sent without an event type')
      return
    }

    server.emit('*', pkg, project)
    server.emit(type, pkg, project)

    reply(res, 200)
    return
  })
}

http.createServer(handler).listen(config.port, config.host, () => {
  log(`Started listening at ${config.host}:${config.port}`)
})

module.exports = {
  http,
  server
}
