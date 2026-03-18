#include "oenwire.h"

/*
 * This code is really messy because it manipulates registers directly.
 * You should not mess with this file unless you really understand what you are
 * doing (i don't lol). I have tried my best to document what my shitty code
 * does here, with a billion comments, but just keep in mind that it is still
 * direct register manipulation so its still weird.
 */

// For reference, we dont set the onewire bus to HIGH in this code ever. instead we set the pin mode to "input" which
// releases the bus high again since its pulled high by a resistor to VDD. we only ever pull it low.

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

// OneWire Basic Functions //
/* ----------------------------- */

// Resets the device on the onewire bus
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

// Writes a bit to the onewire bus
void onewire_write_bit(const uint8_t bit) {
  // if bit = 1
  if (bit) {
    // set the pin as an output
    ow_set_pin_as_output();
    // left-shift the onewire pin bitmask by 16, which resets the pin thus pulling it low
    ONEWIRE_PORT->BSRR = ONEWIRE_PIN << 16;
    // keep it low for 6us - all of these delays are just what the OneWire protocol standard says to use
    delay_us(6);
    // set the pin as an input again
    ow_set_pin_as_input();
    // complete the 70us time slot with the rest 64us
    delay_us(64);
  } else {
    // if the bit is 0 instead, the timeslot is a bit different
    ow_set_pin_as_output();
    ONEWIRE_PORT->BSRR = ONEWIRE_PIN << 16; // we pull the pin low again
    delay_us(60); // this time we keep it low for 60us
    ow_set_pin_as_input(); // set as an input which releases the bus high again
    delay_us(10); // complete the timeslot with 10us more
  }
}

// Reads a bit from the onewire bus (and returns it)
uint8_t onewire_read_bit(void) {
  // again, this is all according to the onewire standard //

  ow_set_pin_as_output(); // set pin as an output
  ONEWIRE_PORT->BSRR = ONEWIRE_PIN << 16; // pull the bus low again
  delay_us(6); // release it for 6us
  ow_set_pin_as_input(); // set the pin as in input, releasing the bus high
  delay_us(9); // wait for another 9us
  // then we do a bitwise AND of the IDR register with the pin's bitmask. this will return a non-zero value if
  // the pin is HIGH. in that case, the ternary operator returns a 1. otherwise, if the pin is low, we return 0.
  // we set that value as the value of the bit variable
  const uint8_t bit = (ONEWIRE_PORT->IDR & ONEWIRE_PIN) ? 1 : 0;
  delay_us(55); // wait another 55us to complete the timeslot
  return bit; // return the bit value
}

// Writes a whole byte to the OneWire bus (8 bits)
void onewire_write_byte(const uint8_t byte) {
  // we iterate from 0 through 7 so we start from the least significant bit (far-right)
  for (int i = 0; i < 8; i++) {
    // we right-shift the byte by the value of i which is the position of the bit in the byte.
    // this, pushes that bit to the last far-right position of the byte. we can then do a bitwise AND
    // with the mask 0x01, so it will either return 1 or 0 depending on whether that bit is set or not in the byte.
    // and we pass that value as the bit argument to the onewire_write_bit() function
    onewire_write_bit((byte >> i) & 0x01);
    // then repeat for all the other bits of the byte
  }
}

// Reads a whole byte from the OneWire bus (8 bits) and returns it
uint8_t onewire_read_byte(void) {
  // we init a uint8_t variable called byte, this will store the byte
  uint8_t byte = 0;

  // we iterate from 0 through 7
  for (int i = 0; i < 8; i++) {
    // the |= operator is the same as saying: byte = byte | [whatever]
    // so we are setting the byte variable to be the value of the byte variable itself, OR'ed with
    // the current state of the bit, left sshifted by i.
    // so if the byte variable was set to 0b00000000 at first, each zero starting from the right, is
    // OR'd with the current bit value, so if its 1 it will set it to 1, eventually completing the byte.
    byte |= (onewire_read_bit() << i);
  }
  return byte; // we return the full byte
}

// OneWire Higher-level Functions //
/* ----------------------------- */

// Begins the conversion
void ds18b20_start_conversion(void) {
  // resets the bus, and returns if it gets a 0 from the reset function which would mean no sensor is present
  if (!onewire_reset()) {
    return;
  }

  // 0xCC skips the ROM address. since we only have one sensor on the bus, we can use this, in order
  // to basically say that we dont care which sensor we are talking to, so any sensor can and
  // should respond (the only one we have)
  onewire_write_byte(0xCC);      // Skip ROM

  // 0x44 tells the sensor to start the conversion.
  // conversion in this context refers to the internal ADC of the DS18B20, converting the analog to a digital value
  onewire_write_byte(0x44);

  // Wait 750ms for the 12-bit resolution conversion to complete
  delay_us(750000);
}

// Reads the temperature value after the conversion is done and returns it as a float
OneWire_Status ds18b20_read_temperature(float *out) {
  if (!onewire_reset()) return OneWire_Error; // return OneWire_Error if no sensor is present
  onewire_write_byte(0xCC); // skip the rom code again, we only have one sensor

  // the command 0xBE reads the scratchpad which is just a tiny memory that holds the last conversion value
  onewire_write_byte(0xBE);

  // the temperature is a 16-bit value, stored as two separate bytes
  // we store the least significant byte first, then the most significant byte
  const uint8_t temp_lsb = onewire_read_byte();
  const uint8_t temp_msb = onewire_read_byte();

  // the scratchpad has 9 bytes in total, we store the first two above, and then
  // we discard the other 7 by reading and not storing them anywhere
  for (int i = 0; i < 7; i++) onewire_read_byte();

  // to combine them, we shift the msb to the left by 8 bits then we OR that value with the
  // lsb. consider that all the bits in the msb were ones, then we have 0b1111111. if that is shifted
  // to the left by 8, then that leaves us with 0b1111111100000000. if we OR tha value
  // with the lsb, then the lsb replaces those zeros at the end, which gives us the two bytes combined.
  const int16_t raw = (temp_msb << 8) | temp_lsb;

  // divide by 16 to get the fractions back
  *out = raw / 16.0;
  return OneWire_OK;
}