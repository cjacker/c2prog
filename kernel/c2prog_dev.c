#include "c2_interface.h"
#include "messages.h"

#include <asm/errno.h>
#include <asm/uaccess.h>
#include <linux/device.h>
#include <linux/fs.h>
#include <linux/gpio.h>
#include <linux/init.h>
#include <linux/module.h>
#include <linux/moduleparam.h>
#include <linux/spinlock.h>

#define DEVICE_NAME "c2prog"
#define CLASS_NAME "c2prog"

static int device_open(struct inode *, struct file *);
static int device_release(struct inode *, struct file *);
static ssize_t device_read(struct file *, char *, size_t, loff_t *);
static ssize_t device_write(struct file *, const char *, size_t, loff_t *);
static int chdev_init(void);
static void chdev_exit(void);
static int c2gpio_init(void);

static int chdev_major;
static struct class *chdev_class = NULL;
static struct device *chdev_dev = NULL;
static char *reply;
static char reply_msg[sizeof(struct message) + 1];
static int dev_open = 0;

static unsigned int c2ck = 0;
static unsigned int c2d = 0;

module_param(c2ck, uint, 0);
MODULE_PARM_DESC(c2ck, "C2CK gpio pin");
module_param(c2d, uint, 0);
MODULE_PARM_DESC(c2d, "C2D gpio pin");

static DEFINE_RAW_SPINLOCK(c2prog_lock);

static int __init c2prog_init(void)
{
	int err;
	printk(KERN_INFO "Registering c2prog module.\n");

	// zero the reply message
	memset(&reply_msg, 0, sizeof(reply_msg));

	// init char device
	err = chdev_init();
	if (err) {
		printk(KERN_ERR "c2prog: Failed to register chardev\n");
		return err;
	}

	// init gpio pins
	err = c2gpio_init();
	if (err) {
		chdev_exit();
		printk(KERN_ERR "c2prog: Failed to setup gpio pins\n");
		return err;
	}

	return 0;
}

static void __exit c2prog_exit(void)
{
	printk(KERN_INFO "Unregistering c2prog module.\n");

	// unregister the character device
	chdev_exit();

	printk(KERN_INFO "c2prog: removing module\n");
}

module_init(c2prog_init);
module_exit(c2prog_exit);

static struct file_operations fops = {
    .read = device_read,
    .write = device_write,
    .open = device_open,
    .release = device_release,
};

static void chdev_exit(void)
{
	device_destroy(chdev_class, MKDEV(chdev_major, 0));
	class_unregister(chdev_class);
	class_destroy(chdev_class);
	unregister_chrdev(chdev_major, DEVICE_NAME);
}

static int chdev_init(void)
{
	// allocate a major number for the device
	chdev_major = register_chrdev(0, DEVICE_NAME, &fops);

	if (chdev_major < 0) {
		printk(KERN_ERR
		       "Registering the c2prog char dev failed with: %d\n",
		       chdev_major);
		return chdev_major;
	}

	chdev_class = class_create(THIS_MODULE, CLASS_NAME);
	if (IS_ERR(chdev_class)) {
		unregister_chrdev(chdev_major, DEVICE_NAME);
		printk(KERN_ERR "Registering the c2prog class failed\n");
		return PTR_ERR(chdev_class);
	}

	chdev_dev = device_create(chdev_class, NULL, MKDEV(chdev_major, 0),
				  NULL, DEVICE_NAME);
	if (IS_ERR(chdev_dev)) {
		class_destroy(chdev_class);
		unregister_chrdev(chdev_major, DEVICE_NAME);
		printk(KERN_ERR "Registering the c2prog device failed\n");
		return PTR_ERR(chdev_dev);
	}

	printk(KERN_INFO "c2prog: registered character device\n");

	return 0;
}

static int device_open(struct inode *inode, struct file *file)
{
	if (dev_open)
		return -EBUSY;

	dev_open = 1;

	return 0;
}

static int device_release(struct inode *inode, struct file *file)
{
	dev_open--;
	return 0;
}

static ssize_t device_read(struct file *file, char *buff, size_t len,
			   loff_t *off)
{
	int bytes_read = 0;

	// if no reply is available, return 0 signifying the end of file
	if (*reply == 0)
		return 0;

	while (len && *reply) {
		put_user(*(reply++), buff++);
		len--;
		bytes_read++;
	}

	return bytes_read;
}

/**
 * process_c2_op() - used for timing critical commands
 */
static void process_c2_op(struct message *msg, struct message *res)
{
	unsigned long flags;

	// disable interrupts to meet the maximum reset low period of the EFM8
	raw_spin_lock_irqsave(&c2prog_lock, flags);

	switch (msg->op) {
	case MESSAGE_OP_RESET:
		C2_reset();
		break;

	case MESSAGE_OP_WRITE_AR:
		C2_write_ar(msg->data[0]);
		break;

	case MESSAGE_OP_READ_AR:
		res->data[0] = C2_read_ar();
		break;

	case MESSAGE_OP_WRITE_DR:
		C2_write_dr(msg->data[0]);
		break;

	case MESSAGE_OP_READ_DR:
		res->data[0] = C2_read_dr();
		break;

	case MESSAGE_OP_POLL:
		C2_poll(msg->data[0], msg->data[1]);
		break;

	case MESSAGE_OP_WRITE_SFR:
		C2_write_sfr(msg->data[0], msg->data[1]);
		break;

	case MESSAGE_OP_READ_SFR:
		res->data[0] = C2_read_sfr(msg->data[0]);
		break;

	case MESSAGE_OP_WRITE_CMD:
		C2_write_cmd(msg->data[0]);
		break;

	case MESSAGE_OP_READ_RESPONSE:
		res->data[0] = C2_read_response();
		break;

	case MESSAGE_OP_READ_DATA:
		res->data[0] = C2_read_data();
		break;

	case MESSAGE_OP_HALT:
		C2_halt();
		break;

	case MESSAGE_OP_DEV_ID:
		res->data[0] = C2_get_dev_id();
		break;
	}

	raw_spin_unlock_irqrestore(&c2prog_lock, flags);
}

static void process_op(struct message *msg, struct message *res)
{
	switch (msg->op) {
	case MESSAGE_OP_PING:
		res->op = MESSAGE_OP_PONG;
		res->data[0] = msg->data[0];
		res->data[1] = msg->data[1];
		break;

	default:
		process_c2_op(msg, res);
		break;
	}
}

static ssize_t device_write(struct file *file, const char *buff, size_t len,
			    loff_t *off)
{
	struct message *msg, *rmsg;

	if (len != sizeof(struct message))
		return -EINVAL;

	msg = (struct message *)buff;
	rmsg = (struct message *)&reply_msg;

	process_op(msg, rmsg);

	// set the reply
	reply = (char *)rmsg;

	return len;
}

static int c2gpio_init(void)
{
	int err;

	err = gpio_request(c2d, "C2D");
	if (err) {
		printk(KERN_ERR "c2prog: Failed to request C2D gpio pin\n");
		return err;
	}

	err = gpio_request(c2ck, "C2CK");
	if (err) {
		printk(KERN_ERR "c2prog: Failed to request C2CK gpio pin\n");
		return err;
	}

	return 0;
}

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Alexander Müller");
MODULE_DESCRIPTION(
    "Programmer implementing the C2 protocol for 8051 based CPUs");
MODULE_VERSION("0.1");