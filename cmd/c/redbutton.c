/* $Id$ */

#include <sys/types.h>
#include <sys/wait.h>
#include <sys/stat.h>

#include <err.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>

#include <usb.h>

#define VENDOR		0x1D34
#define PRODUCT		0x000D
#define SET_REPORT	0x09
#define OUTPUT		0x0200
#define IFACE		0

struct usb_device *
findButton(int vendor, int product)
{
	struct usb_bus *bus;
	struct usb_device *dev;

	usb_init();

	usb_find_busses();
	usb_find_devices();

	for (bus = usb_get_busses(); bus; bus = bus->next)
		for (dev = bus->devices; dev; dev = dev->next)
			if (dev->descriptor.idVendor == vendor &&
			    dev->descriptor.idProduct == product)
				return dev;

	return NULL;
}

int
getButton(struct usb_dev_handle *handle, struct usb_device *dev)
{
	int req = USB_ENDPOINT_OUT|USB_TYPE_CLASS|USB_RECIP_INTERFACE;
	struct usb_endpoint_descriptor *ep =
		dev->config->interface->altsetting->endpoint;
	int addr = ep->bEndpointAddress;
	int interval = ep->bInterval;
	char buf[8];

	bzero(buf, sizeof(buf));
	buf[7] = 2;

	usb_control_msg(handle, req, SET_REPORT, OUTPUT, IFACE,
		buf, sizeof(buf), interval);

	usb_interrupt_read(handle, addr,
		buf, sizeof(buf), interval);

	usleep(interval);

	return buf[0];
}

void
go(char *script, char *action)
{
	int status;

	switch(fork()) {
	case -1:
		err(1, "fork");
		/* NOTREACHED */
	case 0:
		execl(script, script, action, NULL);
		/* NOTREACHED */
	default:
		wait(&status);	
		if (WEXITSTATUS(status) != 0)
			warnx("child failed");
		break;
	}
}

enum { INIT, ARMED, LAUNCH, LOCKED };

int
main(char argc, char **argv)
{
	char *script = "redbutton.sh";
	struct usb_device *dev;
	struct usb_dev_handle *handle;
	int state = INIT;
	int value, button, lid;
	struct stat scr;

	dev = findButton(VENDOR, PRODUCT);
	if (!dev)
		errx(1, "no button found");

	if (argc == 2)
		script = argv[1];

	if (stat(script, &scr) == -1)
		err(1, "%s", script);

	handle = usb_open(dev);
	if (!handle)
		err(1, "%s", usb_strerror());

	usb_detach_kernel_driver_np(handle, IFACE);

	if (usb_claim_interface(handle, IFACE) < 0)
		err(1, "%s", usb_strerror());

	for (;;) {
		value = getButton(handle, dev);
		if (!value)
			continue;

		button = !(value & 0x01);
		lid = !!(value & 0x02);

		switch (state) {
		case INIT:
			if (lid) {
				go(script, "armed");
				state = ARMED;
			}
			break;
		case ARMED:
			if (!lid) {
				go(script, "reset");
				state = INIT;
			}
			if (button) {
				go(script, "launch");
				state = LAUNCH;
			}
			break;
		case LAUNCH:
			if (!button) {
				go(script, "locked");
				state = LOCKED;
			}
			break;
		case LOCKED:
			if (!lid) {
				go(script, "reset");
				state = INIT;
			}
			break;
		}
	}

	usb_release_interface(handle, IFACE);
	usb_close(handle);
	
	return 0;
}
