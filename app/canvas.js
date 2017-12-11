
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

class Canvas {
  constructor(element) {
    this.cursorX = 16
    this.cursorY = 16
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
    drawCanvas.addEventListener('touchstart', this.processTouchMove.bind(this), false)
    drawCanvas.addEventListener('touchmove', this.processTouchMove.bind(this), false)
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

  processTouchMove(evt) {
    // Get the relative coordinates to the canvas.
    // https://stackoverflow.com/a/33756703/101923
    var rect = evt.target.getBoundingClientRect();
    var x = evt.targetTouches[0].pageX - rect.left;
    var y = evt.targetTouches[0].pageY - rect.top;

    this.cursorX = Math.floor(x / this.zoomLevel)
    this.cursorY = Math.floor(y / this.zoomLevel)

    if (!!this.isDrawing) {
      this.context.fillRect(this.cursorX, this.cursorY, 1, 1)
    }

    // Call preventDefault() to prevent mouse events
    evt.preventDefault();
  }

  processMouseDown(evt) {
    this.moving = true
  }

  processMouseUp(evt) {
    this.moving = false
  }

  processMouseMove(evt) {
    if (!this.moving) {
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