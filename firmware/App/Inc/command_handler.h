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

#endif //ODSCI_COMMAND_HANDLER_H
