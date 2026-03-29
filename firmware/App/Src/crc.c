// Copyright 2026 Stratos Thivaios
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

#include "crc.h"

uint8_t generate_crc(const uint8_t data[9], const uint8_t length) {
  uint8_t crc = 0;

  for (int i = 0; i < length; i++) {
    uint8_t byte = data[i];
    for (int bit = 0; bit < 8; bit++) {
      if (((byte & 0x01) ^ (crc & 0x01)) == 1) {
        crc = (crc >> 1) ^ 0x8C;
      } else {
        crc = crc >> 1;
      };
      byte = byte >> 1;
    }
  }

  return crc;
}