#!/usr/bin/env python
# Copyright 2016 PixelArtVision contributors. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from __future__ import division

import argparse
import colorsys
import os
import sys

from PIL import Image


PALETTES = {
    'pico-8': [
        (  0,   0,   0),  # black
        ( 32,  51, 123),  # dark_blue
        (126,  37,  83),  # dark_purple
        (  0, 144,  61),  # dark_green
        (171,  82,  54),  # brown
        ( 52,  54,  53),  # dark_gray
        (194, 195, 199),  # light_gray
        (255, 241, 232),  # white
        (255,   0,  77),  # red
        (255, 155,   0),  # orange
        (255, 231,  39),  # yellow
        (  0, 226,  50),  # green
        ( 41, 173, 255),  # blue
        (132, 112, 169),  # indigo
        (255, 119, 168),  # pink
        (255, 214, 197),  # peach
    ],
    # From colorwheel palette at:
    # https://forums.tigsource.com/index.php?topic=25396.0
    'neon': [
        (255, 255, 255),
        (0, 0, 0),
        (255, 194, 219),
        (188, 255, 153),
        (0, 255, 65),
        (255, 0, 188),
        (255, 0, 124),
        (255, 0, 60),
        (255, 0, 0),
        (255, 64, 0),
        (255, 128, 0),
        (255, 192, 0),
        (254, 255, 0),
        (190, 255, 0),
        (126, 255, 0),
        (63, 255, 0),
        (0, 255, 1),
    ],
    'gameboy': [
        (15, 56, 15),
        (48, 98, 48),
        (139, 172, 15),
        (155, 188, 15),
    ],
}


def colordist(pix1, pix2):
    pix1 = colorsys.rgb_to_hsv(*[channel / 256 for channel in pix1])
    pix2 = colorsys.rgb_to_hsv(*[channel / 256 for channel in pix2])
    dist = 0
    for i in range(3):
        dist += abs(pix1[i] - pix2[i]) #** 2
    return dist


def getclosest(pixel, palette):
    mindistcolor = palette[0]
    mindist = colordist(pixel, palette[0])
    for color in palette:
        dist = colordist(pixel, color)
        if dist < mindist:
            mindist = dist
            mindistcolor = color
    return mindistcolor


def resize(im, width=128):
    owidth, oheight = im.size
    ratio = oheight / owidth
    height = int(ratio * width)
    return im.resize((width, height), Image.LANCZOS)


def main(in_path, out_path, palette='neon', width=128):
    im = Image.open(in_path)
    im = resize(im, width)
    w, h = im.size

    for x in range(w):
        for y in range(h):
            pix = im.getpixel((x,y))
            im.putpixel(
                (x,y),
                getclosest(
                    pix,
                    palette=PALETTES[palette]))

    im.save(out_path)


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        'input', help='Full path of image file to be quantized')
    parser.add_argument(
        'output', help='Full path of where to write quantized image file')
    parser.add_argument(
        '--width', help='Width of output.', default=128, type=int)
    parser.add_argument(
        '--palette', help='Which color palette to use.', default='neon')

    args = parser.parse_args()

    main(args.input, args.output, args.palette, width=args.width)
