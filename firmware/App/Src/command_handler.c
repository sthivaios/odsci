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

#include "../Inc/command_handler.h"

#include "../Inc/version.h"
#include "onewire.h"
#include "usbd_cdc.h"
#include "usbd_cdc_if.h"

#include <stdio.h>

#include "../Inc/led_control.h"


static char Buffer[BUFFER_SIZE];
uint32_t bufferIndex = 0;

void handle_cmd() {
  char tx_buf[256];

  if (strcasecmp(Buffer, "GET_TEMPERATURE") == 0) {
    const TakeAction_Params_T params = {
      .sendTemperature = true
    };
    change_takeAction_params(params);
  } else if (strcasecmp(Buffer, "GET_INFO") == 0) {
    const TakeAction_Params_T params = {
      .sendInfo = true
    };
    change_takeAction_params(params);
  } else if (strcasecmp(Buffer, "HELLO") == 0) {
    snprintf(tx_buf, sizeof(tx_buf), "================================\r\nHey there!\r\nThis is ODSCI v%s\r\nStratos Thivaios (c) 2026\r\n================================\r\n", FIRMWARE_VERSION_STR);
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  } else if (strcasecmp(Buffer, "SET_CLED_ON") == 0) {
    led_control(CAPTURE_LED, true);
  } else if (strcasecmp(Buffer, "SET_CLED_OFF") == 0) {
    led_control(CAPTURE_LED, false);
  } else if (strcasecmp(Buffer, "PING") == 0) {
    snprintf(tx_buf, sizeof(tx_buf), "Pong!\r\n");
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  } else if (strcmp(Buffer, "") == 0) {
    snprintf(tx_buf, sizeof(tx_buf), "ERROR:NO_COMMAND_ENTERED\r\n");
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  } else {
    snprintf(tx_buf, sizeof(tx_buf), "ERROR:UNKNOWN_COMMAND\r\n");
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  }
}

void odsci_handle_rx(const uint8_t *IncBuf, uint32_t Len) {
  for (int i = 0; i < Len; i++) {
    if (bufferIndex >= BUFFER_SIZE) {
      bufferIndex = 0;
    }
    const char byte = (char)IncBuf[i];
    Buffer[bufferIndex] = byte;
    if (Buffer[bufferIndex] == '\r') {
      Buffer[bufferIndex] = '\0';
      handle_cmd();
      bufferIndex = 0;
      continue;
    }
    bufferIndex++;
  }
};

void take_action(const TakeAction_Params_T params) {
  if (params.sendTemperature == true) {
    static float temperature;
    led_control(ACTIVITY_LED, true);
    ds18b20_start_conversion();
    const OneWire_Status status = ds18b20_read_temperature(&temperature);
    led_control(ACTIVITY_LED, false);

    char tx_buf[128];
    if (status == OneWire_Error) {
      snprintf(tx_buf, sizeof(tx_buf), "ERROR:SENSOR_ERROR\r\n");
      CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
      return;
    }
    snprintf(tx_buf, sizeof(tx_buf), "%f\r\n", temperature);
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  } else if (params.sendInfo == true) {
    char tx_buf[128];
    snprintf(tx_buf, sizeof(tx_buf), "ODSCI v%s\r\n", FIRMWARE_VERSION_STR);
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  }
}