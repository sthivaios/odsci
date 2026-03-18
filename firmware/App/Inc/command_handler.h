#ifndef ODSCI_COMMAND_HANDLER_H
#define ODSCI_COMMAND_HANDLER_H

// includes
#include <stdint.h>
#include "stm32f4xx_hal.h"

// defines
#define BUFFER_SIZE 128

// function declarations
void odsci_handle_rx(const uint8_t *IncBuf, uint32_t Len);

#endif //ODSCI_COMMAND_HANDLER_H
