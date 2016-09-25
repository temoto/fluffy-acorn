Ractive = require('ractive/ractive.runtime')
Spinner = require('spin')

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

module.exports = () =>
  r = new Ractive({
    magic: false
    modifyArrays: false
    template: require('./messages.ract.jade')
    data: {
      messages: [],
    }
  })

  self =
    r: r
    ignoreUsername: 'temoto'
    limit: 13
    lastSeen: null
    updateInterval: 21312
    updateSpinner: newSpinner(null, null)
    updateTimer: null
    updateFetch: null
    render: () => r.render.apply(r, arguments)

    stop: () ->
      self.updateSpinner.stop()
      clearTimeout(self.updateTimer)

    update: () ->
      console.log('update', self)
      self.updateSpinner.spin()
      $('h2', '#recent-messages').append(self.updateSpinner.el)
      url = "#{acorn.apiUrl}/api/recent-messages?last_seen=#{encodeURIComponent(self.lastSeen or '')}&limit=#{self.limit}"
      fetch(url, {method: 'GET'})
        .then((response) ->
          self.updateSpinner.stop()
          if acorn.stateRunning
            self.updateTimer = setTimeout(self.update, self.updateInterval)

          if response.status is 200
            return response.json()
          else
            throw new Error('messages.fetch: ' + response.statusText)
        ).then((result) ->
          if result.retryAfter >= 5
            self.updateInterval = result.retryAfter * 1000

          fetchedMessages = result.data.message_list
          if fetchedMessages.length > 0
            newest = fetchedMessages[0].created_at
            if self.lastSeen
              acorn.notify("Localbitcoins #{m.sender.username}: #{m.msg}") for m in fetchedMessages when (m.created_at > self.lastSeen) and (m.sender.username isnt self.ignoreUsername)
            self.lastSeen = newest

          dataMessages = self.r.get('messages')
          Array.prototype.unshift.apply(dataMessages, fetchedMessages)
          if self.limit > 0
            dataMessages = dataMessages[0..(self.limit-1)]
          self.r.set('messages', dataMessages)
        )

  return self
