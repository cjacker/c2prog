#ifndef MESSAGES_H
#define MESSAGES_H

#include <linux/types.h>

enum message_op {
	MESSAGE_OP_PING = 1,
	MESSAGE_OP_PONG = 2,
	MESSAGE_OP_RESET = 3,
	MESSAGE_OP_WRITE_AR = 4,
	MESSAGE_OP_READ_AR = 5,
	MESSAGE_OP_WRITE_DR = 6,
	MESSAGE_OP_READ_DR = 7,
	MESSAGE_OP_POLL = 8,
	MESSAGE_OP_WRITE_SFR = 9,
	MESSAGE_OP_READ_SFR = 10,
	MESSAGE_OP_WRITE_CMD = 11,
	MESSAGE_OP_READ_RESPONSE = 12,
	MESSAGE_OP_READ_DATA = 13,
	MESSAGE_OP_HALT = 14,
	MESSAGE_OP_DEV_ID = 15
};

struct message {
	uint8_t op;
	uint8_t data[2];
	uint8_t code;
} __attribute__((packed));

#endif
