# ODSCI Firmware

## Supported hardware

### MCU
As of March 31st 2026, the current versions of the firmware are written for the Arm® Cortex-M4 platform, more specifically, the STM32F446RET6, as I am using a NUCLEO-F446RE development board.

The first revision of the hardware (v2.x), uses the STM32F042K6T6, which is based on the Arm® Cortex-M0 architecture.

The firmware will therefore require porting, in order to be able to run on the custom ODSCI hardware.

However, again, as of now, the supported hardware is the STM32F446RET6.

### Sensor (DS18B20)
Aside from the obvious (the microprocessor), the only supported sensor at the moment is the DS18B20 from Analog Devices. You can use any pre-built sensor probe that uses this sensor, such as this one from Adafruit: https://www.adafruit.com/product/381/.

Note that you can extend the cable of course as long as you'd like, since the sensor is digital, so it shouldn't affect connectivity with the ODSCI board, especially since the firmware performs CRC8 data validation to ensure that the data is not corrupted by noise on the data line.

## Flashing instructions
The ODSCI board has a DFU button. Hold this button down while plugging the board in, in order to enter the DFU mode. You can release the button once plugged in.

You can then use any tool that lets you flash a binary via DFU. I recommend `dfu-util`. In the future, if I have enough time, as I'm currently quite busy with school, I might work on a way to update the firmware on the board directly through the ODSCI CLI, so that no external tools are required.

If you do use `dfu-util`, you should use it as follows:

```shell
dfu-util -a 0 -s 0x08000000:leave -D odsci-firmware-v1.1.0-arm-cortex-m4-stm32f446re.bin
```

where, of course, instead of this file, you would have the filename of the binary that you are flashing. You can download the latest firmware binary from the releases page of the repository on GitHub, which is probably where you are reading this right now.

Otherwise, you can, of course, always use tools such as the official **STM32CubeProg from ST**, the programming software from Segger, or anything else.

## Building from source
If you wish to build the firmware from source, clone the repository, and enter the `/firmware` directory.

Then, configure and build the project with CMake as follows:

```shell
cmake -S . --preset Release
cmake --build build --preset Release
```

You can then find the binary or the .hex file in `./build/Release/`.

## Compatibility and changelog

**For the changelog (only for the firmware changelog!), see [CHANGELOG.md](CHANGELOG.md).**

### Compatibility
Nothing to see here yet, I'll add a table soon once new features are added, and the CLI requires a specific firmware version to utilize those features.