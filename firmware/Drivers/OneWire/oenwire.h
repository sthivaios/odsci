#ifndef ODSCI_OENWIRE_H
#define ODSCI_OENWIRE_H

// includes
#include <stdint.h>
#include "stm32f4xx_hal_tim.h"

// defines
#define ONEWIRE_PIN_N 5
#define ONEWIRE_PIN GPIO_PIN_5

// types
typedef enum {
  OneWire_OK = 0,
  OneWire_Error = 1
} OneWire_Status;

#endif //ODSCI_OENWIRE_H
