kd  = require 'kd'
Ace = require 'ace/ace'


module.exports = class IDEAce extends Ace

  focus: ->

    return  if kd.singletons.appManager.frontApp.isChatInputFocused?()

    super
