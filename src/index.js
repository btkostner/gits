/**
 * src/index.js
 * Does the logic behind the hooks
 */

const debug = require('debug')

const config = require('../config')
const projects = require('../projects')
const server = require('./server')

const listen = server.server
const log = debug('gits')

config.projects.forEach((project) => {
  if (project.type == null) {
    console.error(`${project.owner}/${project.repo} has no type value`)
    process.exit(1)
  }

  if (projects[project.type] == null) {
    console.error(`${project.owner}/${project.repo} has an invalid type value`)
    process.exit(1)
  }

  if (project.secret == null) {
    log(`${project.owner}/${project.repo} has no secret. This is insecure`)
  }
})

listen.on('ping', (res, project) => {
  log(`Ping event for ${project.owner}/${project.repo}`)
})

listen.on('create', (res, project) => {
  if (res.ref_type !== 'branch') return
  if (!/^refs\/heads\/.*/.test(res.ref)) return

  const br = res.ref.split('/')[2]
  log(`New branch "${br}" being created for ${project.owner}/${project.repo}`)

  if (projects[project.type]['create'] != null) {
    return projects[project.type]['create'](res, project, br)
    .catch((err) => {
      console.error(`Error while creating ${project.owner}/${project.repo}#${br}`)
      console.error(err)
    })
  }

  log(`No action took for ${project.owner}/${project.repo}#${br} create`)
})

listen.on('push', (res, project) => {
  if (!/^refs\/heads\/.*/.test(res.ref)) return

  const br = res.ref.split('/')[2]
  log(`Push occured for ${project.owner}/${project.repo}#${br}`)

  if (res.created && projects[project.type]['create'] != null) {
    return projects[project.type]['create'](res, project, br)
    .catch((err) => {
      console.error(`Error while creating ${project.owner}/${project.repo}#${br}`)
      console.error(err)
    })
  } else if (res.deleted && projects[project.type]['delete'] != null) {
    return projects[project.type]['delete'](res, project, br)
    .catch((err) => {
      console.error(`Error while deleting ${project.owner}/${project.repo}#${br}`)
      console.error(err)
    })
  } else if (projects[project.type]['push'] != null) {
    return projects[project.type]['push'](res, project, br)
    .catch((err) => {
      console.error(`Error while updating ${project.owner}/${project.repo}#${br}`)
      console.error(err)
    })
  }

  log(`No action took for ${project.owner}/${project.repo}#${br} push`)
})

listen.on('delete', (res, project) => {
  if (!/^refs\/heads\/.*/.test(res.ref)) return

  const br = res.ref.split('/')[2]
  log(`Delete occured for ${project.owner}/${project.repo}#${br}`)

  if (res.deleted && projects[project.type]['delete'] != null) {
    return projects[project.type]['delete'](res, project, br)
    .catch((err) => {
      console.error(`Error while deleting ${project.owner}/${project.repo}#${br}`)
      console.error(err)
    })
  }

  log(`No action took for ${project.owner}/${project.repo}#${br} delete`)
})
