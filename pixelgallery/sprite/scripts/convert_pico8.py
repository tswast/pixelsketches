# Copyright 2018 Tim Swast
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Script to convert 'gamekitty' PICO-8 state machine to JSON."""

import json
import re


PICO8_STATE = """states={
 { --1: sitting
 	s={0}, --sprites for state animation
 	next={1,2,8}, -- next possible states
 },
 { --2: stand-up
  s={1,2,3,4},
  next={3,6,9},
 },
 { --3: stand to walk
  s={34,35},
  next={4},
 },
 { --4: walk
  s={
   36,37,38,39,
   --skip 40, standing
  },
  next={5},
  xo=-1, --move left at end
 },
 { --5: walk, other foot
  s={41,42,43,44},
  next={4,6},
  xo=-1,  --move left at end
 },
 { --6: stand waggle
  s={26,4},
  --next={4,6,7,9},
  next={10},
 },
 { --7: sit down
  s={57,58,59,60,61,62,63,64,65,66},
  next={1},
 },
 { --8: turn from sitting
  s={10,11,12,13,14,15,16},
  next={3},
  flipx=true,
 },
 { --9: turn from standing
  s={18,19,20,21,22,23,24,25},
  next={6},
  flipx=true,
 },
 { --10: stand to jump
  s={46,47,48,49},
  xo=-1,
  next={11},
 },
 { --11: jump, start
  s={50,51,45},
  next={12},
  xo=-1,
 },
 { --12: jump, end
  s={52,48,49},
  next={11},
  xo=-1,
 },
}"""


# Make whole thing into a JSON array.
states_lines = PICO8_STATE.split('\n')
del states_lines[0]
del states_lines[-1]
states_lines.insert(0, '[')
states_lines.append(']')
states = '\n'.join(states_lines)

# Remove comments.
states = re.sub(r'--.*$', '', states, flags=re.MULTILINE)

# Convert to JSON keys.
states = re.sub(r'([a-z]+)=', r'"\1": ', states)

# Convert to JSON arrays.
states = re.sub(r'\{((,|\d|\s|\n)+)\}', r'[\1]', states)

# Remove trailing commas.
states = re.sub(r',(\s|\n)*\}', r'}', states)
states = re.sub(r',(\s|\n)*\]', r']', states)
#print(states)

parsed_states = json.loads(states)
print(json.dumps(parsed_states, indent=2))
