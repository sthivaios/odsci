#include "oenwire.h"

/*
 * This code is really messy because it manipulates registers directly.
 * You should not mess with this file unless you really understand what you are
 * doing (i don't lol). I have tried my best to document what my shitty code
 * does here, with a billion comments, but just keep in mind that it is still
 * direct register manipulation so its still weird.
 */

// Helper Functions //
/* ----------------------------- */

// Sets the onewire pin to output mode by manipulating the MODER register
void ow_set_pin_as_output(void) {
  // below, we left-shift 0b11 by the pin number multiplied by 2, to create a mask that has
  // all 0's except those two pins that are 1's. then we do a bitwise NOT to the mask, so that
  // we invert it. now they are all 1's except those two pins that are now 0's.
  // then we do a bitwise AND with the register value, to clear those two bits in the register
  ONEWIRE_PORT->MODER &= ~(0b11 << (ONEWIRE_PIN_N * 2));

  // then we left shift 0b1 by the pin number multiplied by 2. this means that those
  // two pins that we set to 00 above, are now going to become 01, which is how the pin is set
  // to output mode according to the ref manual
  ONEWIRE_PORT->MODER |=  (0b1 << (ONEWIRE_PIN_N * 2));
}

// Sets the onewire pin to input mode by manipulating the MODER register
void ow_set_pin_as_input(void) {
  // does the same exact thing that the first line in the ow_set_pin_as_output() does
  // except this time we just keep it like this, unlike the function above where
  // we then set the two bits to 01. here we leave them as 00 which is the configuration
  // for the input mode.
  ONEWIRE_PORT->MODER &= ~(0b11 << (ONEWIRE_PIN_N * 2));
}

// Delays by a specific number of microseconds
void delay_us(uint32_t microseconds) {
  // HAL_Delay is too slow for this so we use this function instead

  // we set the "start" variable to the current timer value
  uint32_t start = __HAL_TIM_GET_COUNTER(&htim2);

  // do nothing while the difference between the current and the start value is smaller than
  // the desired delay in microseconds
  while ((__HAL_TIM_GET_COUNTER(&htim2) - start) < microseconds);
}

// OneWire Functions //
/* ----------------------------- */

uint8_t onewire_reset(void) {
  // set the pin to output mode
  ow_set_pin_as_output();

  // according to the reference manual the first 15 bits of the BSSR register correspond
  // to the pins of the GPIO port respectively, and are used for setting a pin, whereas
  // the other 15 from 16 to 31 are used for resetting a pin.

  // here, we left-shift the onewire pin bitmask by 16, which resets the pin thus pulling it low
  ONEWIRE_PORT->BSRR = ONEWIRE_PIN << 16;
  // the OneWire protocol requires the line to be pulled low for 480us to reset the device
  delay_us(480);
  // we set the line as an input again
  ow_set_pin_as_input();
  // the OneWire protocol requires a 70us wait after releasing the line HIGH again, before sampling it
  delay_us(70);

  // we sample the onewire line by doing a bitwise AND with the bitmask of the onewire pin.
  // if the sensor is present it pulls the line low, therefore this will return all zeros
  // otherwise it will return a non-zero value if there is no sensor.
  // so we flip the value and set it as the value of the "presence" variable, where 1 now means a sensor is present
  uint8_t presence = !(ONEWIRE_PORT->IDR & ONEWIRE_PIN);

  // according to the OneWire protocol standard we need to wait 410us more to complete the 960us timeslot
  delay_us(410); // complete timeslot

  // we return the sensor presence value
  return presence;
}