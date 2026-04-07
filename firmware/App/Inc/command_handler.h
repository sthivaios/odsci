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

#ifndef ODSCI_COMMAND_HANDLER_H
#define ODSCI_COMMAND_HANDLER_H

// includes
#include <stdint.h>
#include "stm32f4xx_hal.h"

#include <stdbool.h>

// defines
#define BUFFER_SIZE 128

// typedefs
typedef struct {
  bool sendTemperature;
  bool sendInfo;
} TakeAction_Params_T;

// function declarations
void odsci_handle_rx(const uint8_t *IncBuf, uint32_t Len);
void take_action(const TakeAction_Params_T params);
void set_last_reset_due_to_iwdg(const bool iwdg_status);
void get_serial_number(char *out, size_t len);

#endif //ODSCI_COMMAND_HANDLER_H
