'use strict'
Spinner = require('spin')
Ractive = require('ractive')
window.acorn = window.acorn || {}
(() ->
  # acorn.comp = require('./component/messages.html')

  main = () ->
    if !window.fetch?
      $('#msg-invalid-browser')
        .attr('style', 'display:block')
        .append('<p>require window.fetch</p>')

    acorn.setState('pause')

    $(document.body).append('<div id=ract-test></div>')
    acorn.comprecent = new Ractive({el: '#ract-test', template: require('./component/messages.jade')})

  acorn.apiUrl = 'http://127.0.0.1:8001'
  acorn.notifyEnabled = false

  acorn.notify = (msg) ->
    if acorn.notifyEnabled
      new Notification(msg)

  acorn.stateRunning = false
  acorn.setState = (s) ->
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
        acorn.setState(if acorn.stateRunning then 'pause' else 'run')
      else
        throw ('acorn.setState: invalid value: ' + s)

  if window.Notification?
    switch (Notification.permission)
      when 'granted'
        acorn.notifyEnabled = true
      when 'denied'
      else
        Notification.requestPermission( (permission) ->
          acorn.notifyEnabled = permission is 'granted'
        )

  newSpinner = (options, target) ->
    opts =
      lines: 9
      length: 7
      width: 10
      radius: 14
      scale: 1
      corners: 1
      color: '#000'
      opacity: 0.25
      rotate: 0
      direction: 1
      speed: 1
      trail: 51
      fps: 20
      zIndex: 2e9
      className: 'spinner'
      position: 'absolute'
      top: '2rem'
      left: '50%'
      shadow: false
      hwaccel: false

    spinner = new Spinner(opts)

    if target?
      spinner.spin(target)
    return spinner

  acorn.recentMessages =
    fill: null
    update: null
    timer: null
    spinner: newSpinner(null, null)
    interval: 27300
    lastSeen: null

    stop: () ->
      acorn.recentMessages.spinner.stop()
      clearTimeout(acorn.recentMessages.interval)

    update: () ->
      acorn.recentMessages.spinner.spin()
      $('h2', '#recent-messages').append(acorn.recentMessages.spinner.el)
      $.ajax({
        type: 'GET',
        url: acorn.apiUrl + '/api/recent-messages',
        dataType: 'json',
        timeout: 15000,
        success: (data, status, xhr) ->

          if data.check_interval >= 5
            acorn.recentMessages.interval = data.check_interval * 1000

          if data.message_list.length > 0
            newest = data.message_list[0].created_at
            if acorn.recentMessages.lastSeen
              acorn.notify("Localbitcoins #{m.sender.username}: #{m.msg}") for m in data.message_list when m.created_at > acorn.recentMessages.lastSeen
            acorn.recentMessages.lastSeen = newest

          acorn.recentMessages.fill(data)

        error: (xhr, errorType, error) ->
          console.error('recent-messages:', errorType, error)

        complete: () ->
          acorn.recentMessages.spinner.stop()
          if acorn.stateRunning
            acorn.recentMessages.timer = setTimeout(acorn.recentMessages.update, acorn.recentMessages.interval)
      })

    fill: (data) ->
      tbody = $('#recent-messages-list')
      $('tr', tbody).remove()
      fill = (m) -> $("<tr><td><a href='https://localbitcoins.com/request/online_sell_buyer/#{m.contact_id}'>#{m.contact_id}</a></td><td>#{m.sender.name}</td><td>#{m.msg}</td><td>#{m.created_at}</td></tr>").appendTo(tbody)
      fill m for m in data.message_list[0..9]

  $(document).ready(main)
  $(document).on('click', '#btn-pause', (e) ->
    acorn.setState('toggle')
  )

)();
