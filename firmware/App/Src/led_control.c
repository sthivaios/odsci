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

#include "../Inc/led_control.h"

void led_control(const uint64_t led_pin, const bool state) {
  if (!state) /* if state is false */ {
    // reset the bit to normal state which means the pin will be low as configured in MX
    GPIOA->BSRR = led_pin << 16;
  } else {
    // set the bit so pin goes high
    GPIOA->BSRR = led_pin;
  }
}