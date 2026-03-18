// Copyright (c) 2026 Stratos Thivaios
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef ODSCI_ONEWIRE_H
#define ODSCI_ONEWIRE_H

// includes
#include "stm32f4xx_hal.h"
#include "main.h"

// defines
#define ONEWIRE_PIN_N 5
#define ONEWIRE_PIN GPIO_PIN_5
#define ONEWIRE_PORT GPIOB

// types
typedef enum {
  OneWire_OK = 0,
  OneWire_Error = 1
} OneWire_Status;

// function declarations
void ds18b20_start_conversion(void);
OneWire_Status ds18b20_read_temperature(float *out);

#endif //ODSCI_ONEWIRE_H
