'use strict'
window.acorn = window.acorn || {}
(() =>
  acorn.help = require('./component/help.coffee')()
  acorn.recentMessages = require('./component/messages.coffee')()

  main = () =>
    if !window.fetch?
      $('#msg-invalid-browser')
        .attr('style', 'display:block')
        .append('<p>require window.fetch</p>')

    acorn.setState('pause')
    acorn.recentMessages.render('#recent-messages-cont')

    acorn.setupNotification()

    acorn.setKeyScope('main')

  acorn.apiUrl = 'http://127.0.0.1:8001'
  acorn.notifyEnabled = false

  acorn.notify = (msg) =>
    if acorn.notifyEnabled
      new Notification(msg)

  acorn.stateRunning = false
  acorn.setState = (s) =>
    switch (s)
      when 'run'
        acorn.stateRunning = true
        acorn.recentMessages.update()
        $('#app-state').removeClass('app-state-paused')
        $('#app-state').addClass('app-state-running')
        $('#btn-pause').text('running')
      when 'pause'
        acorn.stateRunning = false
        acorn.recentMessages.stop()
        $('#app-state').addClass('app-state-paused')
        $('#app-state').removeClass('app-state-running')
        $('#btn-pause').text('paused')
      when 'toggle'
        acorn.toggleState()
      else
        throw ('acorn.setState: invalid value: ' + s)
  acorn.toggleState = () => acorn.setState(if acorn.stateRunning then 'pause' else 'run')

  acorn.setKeyScope = (scope) =>
    Mousetrap.reset()
    switch (scope)
      when 'main'
        Mousetrap.bind(['p', 'з'], acorn.toggleState)
        Mousetrap.bind(['u', 'г'], acorn.update)
        Mousetrap.bind(['?', ','], acorn.help.show)

  acorn.update = () =>
    acorn.recentMessages.update()

  acorn.setupNotification = () =>
    if !window.Notification?
      return
    switch (Notification.permission)
      when 'granted'
        acorn.notifyEnabled = true
      when 'denied'
      else
        Notification.requestPermission( (permission) =>
          acorn.notifyEnabled = permission is 'granted'
        )

  $(document).ready(main)
  $(document).on('click', '#btn-pause', acorn.toggleState)
)()
