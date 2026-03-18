#include "../Inc/command_handler.h"

#include "usbd_cdc.h"
#include "usbd_cdc_if.h"

#include <stdio.h>


static char Buffer[BUFFER_SIZE];
uint32_t bufferIndex = 0;

void handle_cmd() {
  char tx_buf[BUFFER_SIZE+2];
  snprintf(tx_buf, sizeof(tx_buf), "got: %s\r\n", Buffer);
  CDC_Transmit_FS((uint8_t *)tx_buf, strlen(tx_buf));
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