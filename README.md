# c2prog - Programmer for EFM8 C2 interface using the Linux GPIO API

c2prog is a linux programmer for EFM8 chips from Silicon Labs using linux GPIO pins.

I use a Raspberry PI as programmer, but other boards or other linux systems with fast GPIO pins may also work.

Currently only the following chips are supported:
- EFM8SB1


To implement new chips, please see the *Hacking* section.

## Installation
1. Kernel module:

**Please make sure that you have installed your kernel sources before building the module.**

Enter the following commands to build and install the kernel module:
```sh
cd kernel
make
make install
modprobe c2prog c2d=23 c2ck=24
```
In this example a Raspberry PI was used and an EFM8 was connected to GPIO pins 23 and 24.
If your setup differs, please change the pins in the modprobe command.

If you want to reload the module with different pins use:
```sh
rmmod c2prog
modprobe c2prog c2d=23 c2ck=24
```

2. Install the userspace utility c2prog:
```
# cd back from kernel sources
cd ..
make
make install

# Test the communication with the kernel module
c2prog reset
```

## Programming chips
Use the c2prog util.

### Chip info
```
c2prog info
```

### Reading Flash contents
```
c2prog read --fw out.bin
```

### Writing Flash contents
```
c2prog flash --fw firmware.bin
```

This will automatically erase the chip, if you want to disable this use:
```
c2prog flash --fw firmware.bin --no-auto-erase
```

After flashing the flash contents will be verified.

### Erasing Chip
```
c2prog erase
```


### Resetting Chip
```
c2prog reset
```

## Hacking

To implement new chips, search for the Application Note of the C2 interface. It specifies the required init sequences for the different EFM8 variants.

Then implement them in the *programmer/programmer.go* file.

Pull-requests are very welcome. There will be an issue used to track new chips. If you need to have a chip implemented, feel free to ask there and I will eventually provide support.

## License
GPL v2
