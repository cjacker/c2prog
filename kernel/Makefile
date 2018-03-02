c2prog-objs := c2_interface.o c2prog_dev.o
obj-m += c2prog.o

all:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules
clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
install:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules_install
