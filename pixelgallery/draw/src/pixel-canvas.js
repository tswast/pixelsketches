
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

const QUADRANTS = [
  {
    lowerBound: -Math.PI,
    upperBound: -Math.PI / 2,
    directions: [
      'left',
      'up-left',  // Double the diagonals because the cardinal directions
      'up-left',  // appear in two quadrants, but diagonals in only one.
      'up'
    ]
  },
  {
    lowerBound: -Math.PI / 2,
    upperBound: 0,
    directions: [
      'up',
      'up-right',
      'up-right',
      'right'
    ]
  },
  {
    lowerBound: 0,
    upperBound: Math.PI / 2,
    directions: [
      'right',
      'down-right',
      'down-right',
      'down'
    ]
  },
  {
    lowerBound: Math.PI / 2,
    upperBound: Math.PI,
    directions: [
      'down',
      'down-left',
      'down-left',
      'left'
    ]
  }
]

const DIRECTIONS = {
  'left': {x: -1, y: 0},
  'down-left': {x: -1, y: 1},
  'down': {x: 0, y: 1},
  'down-right': {x: 1, y: 1},
  'right': {x: 1, y: 0},
  'up-right': {x: 1, y: -1},
  'up': {x: 0, y: -1},
  'up-left': {x: -1, y: -1}
}

function calculateDirection (startX, startY, x, y) {
  var theta = Math.atan2(y - startY, x - startX)
  var quad = 0
  for (var i = 0, len = QUADRANTS.length; i < len; i++) {
    var upperBound = QUADRANTS[i].upperBound
    if (theta < upperBound) {
      quad = i
      break
    }
  }
  var lowerBound = QUADRANTS[quad].lowerBound
  var quadAngle = theta - lowerBound
  var quadDirection = Math.floor(Math.min(3, quadAngle * 4 / (Math.PI / 2)))
  var direction = QUADRANTS[quad].directions[quadDirection]
  return direction
}

class PixelCanvas {
  constructor (element) {
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
    ctx.fillStyle = 'rgb(200, 0, 0)'
    ctx.fillRect(0, 0, 20, 20)
    ctx.fillStyle = 'rgba(0, 0, 200, 0.5)'
    ctx.fillRect(12, 12, 32, 32)
    this.context = ctx

    element.appendChild(canvas)
    element.appendChild(document.createElement('br'))

    var drawCanvas = document.createElement('canvas')
    drawCanvas.classList.add('app-canvas')
    drawCanvas.setAttribute('width', 32 * this.zoomLevel)
    drawCanvas.setAttribute('height', 32 * this.zoomLevel)
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

  startDraw (evt) {
    this.isDrawing = true
    this.context.fillRect(this.cursorX, this.cursorY, 1, 1)
  }

  endDraw (evt) {
    this.isDrawing = false
  }

  processMoveStartTouch (evt) {
    this.isMoving = true
    this.moveStartX = evt.targetTouches[0].screenX
    this.moveStartY = evt.targetTouches[0].screenY
  }

  squaredDistanceToMoveStart (x, y) {
    return Math.pow(x - this.moveStartX, 2) + Math.pow(y - this.moveStartY, 2)
  }

  processMove (x, y) {
    var dist = this.squaredDistanceToMoveStart(x, y)
    if (dist > this.moveRadiusSqr) {
      var dIndex = calculateDirection(this.moveStartX, this.moveStartY, x, y)
      var direction = DIRECTIONS[dIndex]
      this.cursorX = this.cursorX + direction.x
      this.cursorY = this.cursorY + direction.y
      this.moveStartX = x
      this.moveStartY = y

      if (this.isDrawing) {
        this.context.fillRect(this.cursorX, this.cursorY, 1, 1)
      }
    }
  }

  processMoveTouch (evt) {
    // Call preventDefault() to prevent mouse events
    evt.preventDefault()

    var x = evt.targetTouches[0].screenX
    var y = evt.targetTouches[0].screenY
    this.processMove(x, y)
  }

  processMouseDown (evt) {
    this.isMoving = true
    this.moveStartX = evt.offsetX
    this.moveStartY = evt.offsetY
  }

  processMouseUp (evt) {
    this.isMoving = false
  }

  processMouseMove (evt) {
    if (!this.isMoving) {
      return
    }
    var x = evt.offsetX
    var y = evt.offsetY
    this.processMove(x, y)
  }

  render (timestamp) {
    this.drawContext.clearRect(0, 0, this.drawCanvas.width, this.drawCanvas.height)
    this.drawContext.drawImage(this.canvas, 0, 0, this.drawCanvas.width, this.drawCanvas.height)
    this.drawContext.strokeRect(
      this.cursorX * this.zoomLevel,
      this.cursorY * this.zoomLevel,
      this.zoomLevel,
      this.zoomLevel)
  }
}

export { PixelCanvas as default }
