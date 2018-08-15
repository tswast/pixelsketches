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

import os, sys
from PIL import Image
from PIL import ImagePalette

im = Image.open("/Users/me/Desktop/input.png")
im = im.convert("RGB")

pim = Image.new("P", (1, 1), 0)
palette = [
      0,   0,   0,  # black
     32,  51, 123,  # dark_blue
    126,  37,  83,  # dark_purple
      0, 144,  61,  # dark_green
    171,  82,  54,  # brown
     52,  54,  53,  # dark_gray
    194, 195, 199,  # light_gray
    255, 241, 232,  # white
    255,   0,  77,  # red
    255, 155,   0,  # orange
    255, 231,  39,  # yellow
      0, 226,  50,  # green
     41, 173, 255,  # blue
    132, 112, 169,  # indigo
    255, 119, 168,  # pink
    255, 214, 197,  # peach
]
palette = [
15, 56, 15,
48, 98, 48,
139, 172, 15,
155, 188, 15,
]
# https://forums.tigsource.com/index.php?topic=25396.0
palette = [
255, 255, 255,
0, 0, 0,
255, 194, 219,
188, 255, 153,
0, 255, 65,
255, 0, 188,
255, 0, 124,
255, 0, 60,
255, 0, 0,
255, 64, 0,
255, 128, 0,
255, 192, 0,
254, 255, 0,
190, 255, 0,
126, 255, 0,
63, 255, 0,
0, 255, 1,
]
palette = palette + [0] * ((3 * 256) - len(palette))

pim.putpalette(palette)
out = im.quantize(palette=pim)

out.save("out.png")
