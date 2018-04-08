(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);var f=new Error("Cannot find module '"+o+"'");throw f.code="MODULE_NOT_FOUND",f}var l=n[o]={exports:{}};t[o][0].call(l.exports,function(e){var n=t[o][1][e];return s(n?n:e)},l,l.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){

// Copyright 2017 Tim Swast
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const DIRECTIONS = [
  {x: -1, y: 0},  // -pi
  {x: -1, y: 1},  // -3/4 pi
  {x: 0, y: 1},  // -1/2 pi
  {x: 1, y: 1},  // -1/4 pi
  {x: 1, y: 0},  // 0 pi
  {x: 1, y: -1},  // 1/4 pi
  {x: 0, y: -1},  // 1/2 pi
  {x: -1, y: -1},  // 3/4 pi
]

class Canvas {
  constructor(element) {
    this.cursorX = 16
    this.cursorY = 16
    this.isMoving = false
    this.moveStartX = 0
    this.moveStartY = 0
    this.moveRadiusSqr = 16 * 16
    this.zoomLevel = 4

    var canvas = document.createElement('canvas')
    canvas.setAttribute('width', 32)
    canvas.setAttribute('height', 32)
    this.canvas = canvas
  
    // Draw something
    // https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Tutorial/Basic_usage
    var ctx = canvas.getContext('2d')
    ctx.fillStyle = 'rgb(200, 0, 0)';
    ctx.fillRect(0, 0, 20, 20);
    ctx.fillStyle = 'rgba(0, 0, 200, 0.5)';
    ctx.fillRect(12, 12, 32, 32);
    this.context = ctx
  
    element.appendChild(canvas)
    element.appendChild(document.createElement('br'))
  
    var drawCanvas = document.createElement('canvas')
    drawCanvas.classList.add('app-canvas')
    drawCanvas.setAttribute('width', 32*this.zoomLevel)
    drawCanvas.setAttribute('height', 32*this.zoomLevel)
    this.drawCanvas = drawCanvas
  
    // Zoom, but don't mess with my pixels
    // https://codepo8.github.io/canvas-images-and-pixels/
    this.drawContext = drawCanvas.getContext('2d')
    this.drawContext.imageSmoothingEnabled = false
    this.drawContext.webkitImageSmoothingEnabled = false
    this.drawContext.msImageSmoothingEnabled = false
  
    element.appendChild(drawCanvas)
    element.appendChild(document.createElement('br'))
    drawCanvas.addEventListener('touchstart', this.processMoveStartTouch.bind(this), false)
    drawCanvas.addEventListener('touchmove', this.processMoveTouch.bind(this), false)
    drawCanvas.addEventListener('mousedown', this.processMouseDown.bind(this), false)
    drawCanvas.addEventListener('mouseup', this.processMouseUp.bind(this), false)
    drawCanvas.addEventListener('mousemove', this.processMouseMove.bind(this), false)

    var button = document.createElement('button')
    button.textContent = 'Dot'
    element.appendChild(button)
    button.addEventListener('touchstart', this.startDraw.bind(this), false)
    button.addEventListener('mousedown', this.startDraw.bind(this), false)
    button.addEventListener('touchend', this.endDraw.bind(this), false)
    button.addEventListener('mouseup', this.endDraw.bind(this), false)
    this.isDrawing = false
  }

  startDraw(evt) {
    this.isDrawing = true
    this.context.fillRect(this.cursorX, this.cursorY, 1, 1)
  }

  endDraw(evt) {
    this.isDrawing = false
  }

  processMoveStartTouch(evt) {
    this.isMoving = true
    this.moveStartX = evt.targetTouches[0].screenX
    this.moveStartY = evt.targetTouches[0].screenY
  }

  squaredDistanceToMoveStart(x, y) {
    return Math.pow(x - this.moveStartX, 2) + Math.pow(y - this.moveStartY, 2)
  }

  processMoveTouch(evt) {
    // Call preventDefault() to prevent mouse events
    evt.preventDefault();

    var x = evt.targetTouches[0].screenX
    var y = evt.targetTouches[0].screenY
    var dist = this.squaredDistanceToMoveStart(x, y)
    if (dist > this.moveRadiusSqr) {
      var theta = Math.atan2(x - this.moveStartX, y - this.moveStartY)
      // TODO: rotate slighly so cardinal directions match up.
      var dIndex = Math.min(7, Math.floor(8 * (theta + Math.PI) / (2 * Math.PI)))
      var direction = DIRECTIONS[dIndex]
      
      this.cursorX = this.cursorX + direction.x
      this.cursorY = this.cursorY + direction.y
      this.moveStartX = x
      this.moveStartY = y

      if (!!this.isDrawing) {
        this.context.fillRect(this.cursorX, this.cursorY, 1, 1)
      }
    }
  }

  processMouseDown(evt) {
    this.isMoving = true
  }

  processMouseUp(evt) {
    this.isMoving = false
  }

  processMouseMove(evt) {
    if (!this.isMoving) {
      return
    }
    this.cursorX = Math.floor(evt.offsetX / this.zoomLevel)
    this.cursorY = Math.floor(evt.offsetY / this.zoomLevel)

    if (!!this.isDrawing) {
      this.context.fillRect(this.cursorX, this.cursorY, 1, 1)
    }
  }

  render(timestamp) {
    this.drawContext.clearRect(0, 0, this.drawCanvas.width, this.drawCanvas.height)
    this.drawContext.drawImage(this.canvas, 0, 0, this.drawCanvas.width, this.drawCanvas.height)
    this.drawContext.strokeRect(
      this.cursorX * this.zoomLevel,
      this.cursorY * this.zoomLevel,
      this.zoomLevel,
      this.zoomLevel)
  }
}

exports.Canvas = Canvas
},{}],2:[function(require,module,exports){

// Copyright 2017 Tim Swast
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

'use strict'

var canvas = require('./canvas.js')

function render (app) {
  var doRender = function (timestamp) {
    app.render(timestamp)
    window.requestAnimationFrame(doRender)
  }
  return doRender
}

window.addEventListener('load', function () {
  var appElement = document.getElementById('app')
  var app = new canvas.Canvas(appElement)
  window.requestAnimationFrame(render(app))
})

},{"./canvas.js":1}]},{},[2]);
