#include "../Inc/command_handler.h"

#include "onewire.h"
#include "usbd_cdc.h"
#include "usbd_cdc_if.h"

#include <stdio.h>


static char Buffer[BUFFER_SIZE];
uint32_t bufferIndex = 0;

void handle_cmd() {
  char tx_buf[128];

  if (strcasecmp(Buffer, "GET_TEMPERATURE") == 0) {
    const TakeAction_Params_T params = {
      .sendTemperature = true
    };
    change_takeAction_params(params);
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
    ds18b20_start_conversion();
    const OneWire_Status status = ds18b20_read_temperature(&temperature);

    char tx_buf[128];
    if (status == OneWire_Error) {
      snprintf(tx_buf, sizeof(tx_buf), "ERROR:SENSOR_ERROR\r\n");
      CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
      return;
    }
    snprintf(tx_buf, sizeof(tx_buf), "%f\r\n", temperature);
    CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
  }
}