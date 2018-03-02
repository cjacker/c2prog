#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/gpio.h>
#include <linux/delay.h>

#include "c2_interface.h"

#define C2CK_H gpio_set_value(c2ck, 1)
#define C2CK_L gpio_set_value(c2ck, 0)

#define C2D_H gpio_set_value(c2d, 1)
#define C2D_L gpio_set_value(c2d, 0)
#define C2D_R gpio_get_value(c2d)

#define STROBE_C2CK \
    C2CK_H;         \
    C2CK_L;         \
    C2CK_H;

static unsigned int c2d;
static unsigned int c2ck;

static void c2d_input(void)
{
    gpio_direction_input(c2d);
}

static void c2d_output(void)
{
    gpio_direction_output(c2d, 1);
}

static void c2ck_input(void)
{
    //gpio_direction_input(c2ck);
}

static void c2ck_output(void)
{
    gpio_direction_output(c2ck, 1);
}

void C2_init(unsigned int pc2d, unsigned int pc2ck)
{
    c2d = pc2d;
    c2ck = pc2ck;
    c2d_input();
}

void C2_reset(void)
{
    c2ck_output();
    C2CK_L;
    udelay(30);
    C2CK_H;
    udelay(10);
    c2ck_input();
    mdelay(30);
}

void C2_write_ar(uint8_t addr)
{
    uint8_t i;

    c2ck_output();
    STROBE_C2CK;

    c2d_output();

    C2D_H;
    STROBE_C2CK;
    C2D_H;
    STROBE_C2CK;

    for (i = 0; i < 8; i++)
    {
        if (addr & 0x01)
        {
            C2D_H;
        }
        else
        {
            C2D_L;
        }

        STROBE_C2CK;
        addr >>= 1;
    }

    c2d_input();
    STROBE_C2CK;

    c2ck_input();
}

uint8_t C2_read_ar(void)
{
    uint8_t i, addr;

    c2ck_output();

    // START field
    STROBE_C2CK;

    // INS field
    c2d_output();
    C2D_L;
    STROBE_C2CK;
    C2D_H;
    STROBE_C2CK;
    c2d_input();

    addr = 0;
    for (i = 0; i < 8; i++)
    {
        addr >>= 1;
        STROBE_C2CK;

        if (C2D_R == 1)
        {
            addr |= 0x80;
        }
    }

    // STOP field
    STROBE_C2CK;
    c2ck_input();

    return addr;
}

void C2_write_dr(uint8_t dat)
{
    uint8_t i;

    c2ck_output();
    STROBE_C2CK;

    c2d_output();
    C2D_H;
    STROBE_C2CK;
    C2D_L;
    STROBE_C2CK;

    C2D_L;
    STROBE_C2CK;
    C2D_L;
    STROBE_C2CK;

    for (i = 0; i < 8; i++)
    {
        if (dat & 0x01)
        {
            C2D_H;
        }
        else
        {
            C2D_L;
        }

        STROBE_C2CK;
        dat >>= 1;
    }

    c2d_input();
    STROBE_C2CK;

    while (C2D_R == 0)
    {
        STROBE_C2CK;
    }

    STROBE_C2CK;
    c2ck_input();
}

uint8_t C2_read_dr(void)
{
    uint8_t i, dat;
    c2ck_output();
    STROBE_C2CK;

    c2d_output();
    C2D_L;
    STROBE_C2CK;
    C2D_L;
    STROBE_C2CK;

    C2D_L;
    STROBE_C2CK;
    C2D_L;
    STROBE_C2CK;

    c2d_input();

    STROBE_C2CK;

    while (C2D_R == 0)
    {
        STROBE_C2CK;
    }

    dat = 0;
    for (i = 0; i < 8; i++)
    {
        dat >>= 1;
        STROBE_C2CK;
        if (C2D_R == 1)
        {
            dat |= 0x80;
        }
    }

    STROBE_C2CK;
    c2ck_input();

    return dat;
}

void C2_poll(uint8_t flag, uint8_t check)
{
    while (1)
    {
        uint8_t ar = C2_read_ar();

        if (check == 1)
        {
            if ((ar & flag))
            {
                break;
            }
        }
        else
        {
            if ((ar & flag) == 0)
            {
                break;
            }
        }
    }
}

void C2_halt(void)
{
    C2_reset();

    // standard reset
    C2_write_ar(C2_FPCTL);
    C2_write_dr(C2_FPCTL_RESET);

    // core reset
    C2_write_dr(C2_FPCTL_CORE_RESET);

    // HALT
    C2_write_dr(C2_FPCTL_HALT);

    //HAL_Delay(20);
}

uint8_t C2_get_dev_id(void)
{
    C2_reset();
    return C2_read_dr();
}

uint8_t C2_read_sfr(uint8_t addr)
{
    C2_write_ar(addr);
    return C2_read_dr();
}

void C2_write_sfr(uint8_t addr, uint8_t data)
{
    C2_write_ar(addr);
    C2_write_dr(data);
}

void C2_write_cmd(uint8_t cmd)
{
    C2_write_dr(cmd);
    C2_poll(C2_AR_INBUSY, 0);
}

uint8_t C2_read_response(void)
{
    C2_poll(C2_AR_OUTREADY, 1);
    return C2_read_dr();
}

uint8_t C2_read_data(void)
{
    C2_poll(C2_AR_OUTREADY, 1);
    return C2_read_dr();
}