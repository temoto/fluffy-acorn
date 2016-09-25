picoModal = require('picoModal')

module.exports = () =>
  pico = picoModal({
    content: require('pug-html!./help.jade')
    closeHtml: '<span>x</span>'
    closeStyles: {
      position: 'absolute'
      top: '0'
      right: '0'
      background: '#eee'
      padding: '0.3rem 0.6rem'
      cursor: 'pointer'
      borderRadius: '0.3rem'
      border: '.1rem solid #ccc'
    }
  })

  self =
    afterClose: () => acorn.setKeyScope('main')

    close: () => self.close()

    show: () =>
      acorn.setKeyScope('help')
      pico.afterClose(self.afterClose).show()

  return self
