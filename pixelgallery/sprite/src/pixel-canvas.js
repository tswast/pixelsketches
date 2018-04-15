
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

import { default as stateDefinitions } from './state.json'

class PixelCanvas {
  constructor (element) {
    this.slowness = 3
    this.slowCount = 0
    this.zoomLevel = 4
    this.stateIndex = 0
    this.animIndex = 0        
    var state = stateDefinitions[this.stateIndex]
    var frame = state['s'][this.animIndex]
    this.spriteX = (frame * 8) % 128
    this.spriteY = Math.floor(frame / 16) * 8

    var canvas = document.createElement('canvas')
    canvas.setAttribute('width', 32)
    canvas.setAttribute('height', 32)
    this.canvas = canvas

    // Draw something
    // https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API/Tutorial/Basic_usage
    var ctx = canvas.getContext('2d')
    this.spriteSheet = new Image()
    this.spriteSheet.src = './gamekitty_pico8.png'
    ctx.drawImage(this.spriteSheet, this.spriteX, this.spriteY, 8, 8, 2, 10, 8, 8)
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
  }

  update (timeDiff) {
    this.slowCount = this.slowCount + 1
    if (this.slowCount < this.slowness) {
      return
    }
    this.slowCount = 0
    var state = stateDefinitions[this.stateIndex]

    // Move to next frame
    this.animIndex = this.animIndex + 1

    // Move to next state
    if (this.animIndex >= state['s'].length) {
      this.animIndex = 0
      var nextState = Math.floor(Math.random() * state['next'].length)
      this.stateIndex = state['next'][nextState]
      state = stateDefinitions[this.stateIndex]
    }
    var frame = state['s'][this.animIndex]
    this.spriteX = (frame * 8) % 128
    this.spriteY = Math.floor(frame / 16) * 8
  }

  render () {
    this.context.drawImage(this.spriteSheet, this.spriteX, this.spriteY, 8, 8, 2, 10, 8, 8)
    // Draw big version.
    this.drawContext.clearRect(0, 0, this.drawCanvas.width, this.drawCanvas.height)
    this.drawContext.drawImage(this.canvas, 0, 0, this.drawCanvas.width, this.drawCanvas.height)
  }
}

export { PixelCanvas as default }
