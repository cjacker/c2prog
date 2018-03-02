#ifndef C2_INTERFACE_H
#define C2_INTERFACE_H

#include <linux/types.h>

#define C2_AR_OUTREADY 0x01
#define C2_AR_INBUSY 0x02
#define C2_AR_OTPBUSY 0x80
#define C2_AR_OTPERROR 0x40
#define C2_AR_FLBUSY 0x80

#define C2_DEVICEID 0x00
#define C2_REVID 0x01
#define C2_FPCTL 0x02

#define C2_FPCTL_RUNNING 0x00
#define C2_FPCTL_HALT 0x01
#define C2_FPCTL_RESET 0x02
#define C2_FPCTL_CORE_RESET 0x04

//#define C2_FPDAT                0xB4
#define C2_FPDAT_GET_VERSION 0x01
#define C2_FPDAT_GET_DERIVATIVE 0x02
#define C2_FPDAT_DEVICE_ERASE 0x03
#define C2_FPDAT_GET_CRC 0x05
#define C2_FPDAT_BLOCK_READ 0x06
#define C2_FPDAT_BLOCK_WRITE 0x07
#define C2_FPDAT_PAGE_ERASE 0x08
#define C2_FPDAT_DIRECT_READ 0x09
#define C2_FPDAT_DIRECT_WRITE 0x0a
#define C2_FPDAT_INDIRECT_READ 0x0b
#define C2_FPDAT_INDIRECT_WRITE 0x0c

#define C2_FPDAT_RETURN_INVALID_COMMAND 0x00
#define C2_FPDAT_RETURN_COMMAND_FAILED 0x02
#define C2_FPDAT_RETURN_FLASH_ERR 0x03
#define C2_FPDAT_RETURN_NO_HALT_ERR 0x0B
#define C2_FPDAT_RETURN_COMMAND_OK 0x0D

#define C2_DEVCTL 0x02
#define C2_EPCTL 0xDF
#define C2_EPDAT 0xBF
#define C2_EPADDRH 0xAF
#define C2_EPADDRL 0xAE
#define C2_EPSTAT 0xB7
#define C2_EPSTAT_WRITE_LOCK 0x80
#define C2_EPSTAT_READ_LOCK 0x40
#define C2_EPSTAT_CAL_VALID 0x20
#define C2_EPSTAT_CAL_DONE 0x10
#define C2_EPSTAT_ERROR 0x01

#define C2_EPCTL_READ 0x00
#define C2_EPCTL_WRITE1 0x40
#define C2_EPCTL_WRITE2 0x58
#define C2_EPCTL_FAST_WRITE 0x78

void C2_init(unsigned int pc2d, unsigned int pc2ck);
void C2_reset(void);
void C2_write_ar(uint8_t addr);
uint8_t C2_read_ar(void);
void C2_write_dr(uint8_t dat);
uint8_t C2_read_dr(void);
void C2_poll(uint8_t flag, uint8_t check);
void C2_halt(void);
uint8_t C2_get_dev_id(void);
uint8_t C2_read_sfr(uint8_t addr);
void C2_write_sfr(uint8_t addr, uint8_t data);
void C2_write_cmd(uint8_t cmd);
uint8_t C2_read_response(void);
uint8_t C2_read_data(void);

#endif