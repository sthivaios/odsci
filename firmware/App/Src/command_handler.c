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
bool iwdg_reset = false;

void set_last_reset_due_to_iwdg(const bool iwdg_status) {
  iwdg_reset = iwdg_status;
}

void get_serial_number(char *out, size_t len) {
  uint32_t uid0 = *(uint32_t *)(UID_BASE);
  uint32_t uid1 = *(uint32_t *)(UID_BASE + 4);
  uint32_t uid2 = *(uint32_t *)(UID_BASE + 8);

  // mix all three words so the result reflects the full UID
  uint32_t mixed = uid0 ^ uid1 ^ uid2;
  snprintf(out, len, "%08lX", mixed);
}

void handle_cmd() {
  char tx_buf[256];
  ERROR_LED_OFF();

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
  }
#if CLED_IS_FOR_ERRORS_INSTEAD == 0
  else if (strcasecmp(Buffer, "SET_CLED_ON") == 0) {
    led_control(CAPTURE_LED, true);
  } else if (strcasecmp(Buffer, "SET_CLED_OFF") == 0) {
    led_control(CAPTURE_LED, false);
  }
#endif
  else if (strcasecmp(Buffer, "PING") == 0) {
    snprintf(tx_buf, sizeof(tx_buf), "Pong!\r\n");
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  } else if (strcmp(Buffer, "") == 0) {
    ERROR_LED_ON();
    snprintf(tx_buf, sizeof(tx_buf), "ERROR:NO_COMMAND_ENTERED\r\n");
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  } else {
    ERROR_LED_ON();
    snprintf(tx_buf, sizeof(tx_buf), "ERROR:UNKNOWN_COMMAND\r\n");
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  }
}

void odsci_handle_rx(const uint8_t *IncBuf, uint32_t Len) {
  for (int i = 0; i < Len; i++) {
    const char byte = (char)IncBuf[i];
    if (byte == '\r') {
      Buffer[bufferIndex] = '\0';
      handle_cmd();
      bufferIndex = 0;
      continue;
    }
    if (bufferIndex >= BUFFER_SIZE - 1) {
      bufferIndex = 0;
      ERROR_LED_ON();
      static const char err[] = "ERROR:BUFFER_OVERFLOW_COMMAND_TOO_LONG\r\n";
      CDC_Transmit_FS((uint8_t *)err, strlen(err));
      continue;
    }
    Buffer[bufferIndex] = byte;
    bufferIndex++;
  }
};

void take_action(const TakeAction_Params_T params) {
  static char tx_buf[192];
  if (params.sendTemperature == true) {
    OneWire_Status status;

    static float temperature;
    led_control(ACTIVITY_LED, true);
    const OneWire_Status conversion_status = ds18b20_start_conversion();
    if (conversion_status == OneWire_Error) {
      status = OneWire_Error;
    } else {
      status = ds18b20_read_temperature(&temperature);
      led_control(ACTIVITY_LED, false);
    }
    if (status == OneWire_Error) {
      ERROR_LED_ON();
      snprintf(tx_buf, sizeof(tx_buf), "ERROR:SENSOR_ERROR\r\n");
      CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
      return;
    }
    if (status == OneWire_CRC_Error) {
      ERROR_LED_ON();
      snprintf(tx_buf, sizeof(tx_buf), "ERROR:DATA_CRC_ERROR\r\n");
      CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
      return;
    }
    int temp_int = (int)(temperature * 100); // 2 decimal places
    snprintf(tx_buf, sizeof(tx_buf), "%d.%02d\r\n", temp_int / 100, temp_int % 100);
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
    ERROR_LED_OFF();
  } else if (params.sendInfo == true) {
    char serial[24];
    get_serial_number(serial, 24);
    snprintf(tx_buf, sizeof(tx_buf), "FIRMWARE_VERSION=%s,CLED_IS_FOR_ERRORS_INSTEAD=%d,LAST_RESET_DUE_TO_IWDG=%d,SERIAL_NUMBER=%s\r\n", FIRMWARE_VERSION_STR, CLED_IS_FOR_ERRORS_INSTEAD, iwdg_reset, serial);
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  }
}