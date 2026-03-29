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

#ifndef ODSCI_LED_CONTROL_H
#define ODSCI_LED_CONTROL_H

// include
#include <stdbool.h>
#include <stdint.h>

#include "stm32f446xx.h"
#include "stm32f4xx_hal.h"

// led definitions
#define ACTIVITY_LED GPIO_PIN_6
#define CAPTURE_LED GPIO_PIN_7

// function declarations
void led_control(uint64_t led_pin, bool state);

// macros/functions for error indicator
#if CLED_IS_FOR_ERRORS_INSTEAD == 1
#define ERROR_LED_ON()  errorled_on()
#define ERROR_LED_OFF() errorled_off()
void errorled_on(void);
void errorled_off(void);
#else
#define ERROR_LED_ON()  do {} while(0)
#define ERROR_LED_OFF() do {} while(0)
#endif

#endif //ODSCI_LED_CONTROL_H
