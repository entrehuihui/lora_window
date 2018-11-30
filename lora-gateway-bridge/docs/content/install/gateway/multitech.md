---
title: Multitech
menu:
    main:
        parent: gateway
description: Installation of the LoRa Gateway Bridge service on the Multitech Conduit gateway.
---

# Multitech

## Multitech Conduit

* [Product detail page](https://www.multitech.com/brands/multiconnect-conduit)
* [Product documentation page](http://www.multitech.net/developer/products/multiconnect-conduit-platform/)

After completing this steps, you have a Multitech Conduit running both the
packet-forwarder and LoRa Gateway bridge. The packet-forwarder will forwards
the UDP data to `localhost:1700` and the LoRa Gateway Bridge will forward
this data as JSON over MQTT to a MQTT broker. See below:

![Gateway Pieces and Connections](/lora-gateway-bridge/img/MultitechGatewaySettings.png)

There are two different Multitech Conduit firmware versions: mLinux and AEP.
The AEP version comes with a web-interface and IBM Node-RED pre-installed.
The mLinux version provides an open Linux development environment and is
recommended when complete (firmware) control is preferred.

If you don't know the firmware version, here are a couple of ways to tell
them apart:

1. When logging in via the serial port behind the Multitech logo cover, they
   display the type of box they are.
2. The AEP model ships with the default login/password as admin/admin.  The 
   mLinux version uses root/root.

Please refer to [http://www.multitech.net/developer/products/multiconnect-conduit-platform/](http://www.multitech.net/developer/products/multiconnect-conduit-platform/)
for more documentation on on the Multitech Conduit.

**Note:** It is possible to install mLinux on an AEP Conduit version following the
steps below. This is recommended when you don't rely on any software provided by
the AEP firmware.

### Getting the IP address

Before continuing, you'll want to obtain the IP address of the Conduit.  This can 
be done using a serial connection from a computer using a USB-to-microUSB cable,
connecting to the plug behind the Multitech logo placard.  Plug the device into 
your network, provide power, and let it boot until the "STATUS" light is 
blinking in a heartbeat pattern.  Connect to the device via a serial terminal 
program.  Once logged in, issue the command "ifconfig" to get the IP address of 
the eth0 connection.  Note that if the IP address is `192.168.2.1`, the device is
likely configured with a static IP.  In this case, edit the file 
`/etc/network/interfaces`, change the line that says, `iface eth0 inet static` to 
`iface eth0 inet dhcp`, and comment out the lines specifying the IP Address and 
netmask by adding a `#` at the beginning of each line:

{{<highlight text>}}
# address 192.168.2.1
# netmask 255.255.255.0
{{< /highlight >}}

Then execute `/etc/init.d/networking restart`, and obtain the issued IP address
as outlined above.

### Upgrading / migrating from AEP to the latest mLinux

The suggested way to setup the packet-forwarder and LoRa Gateway Bridge on a
Multitech Conduit is by using the base mLinux firmware image
(`mlinux-base-*.jffs2`). This firmware
image installs the minimal amount of software needed to boot the Conduit,
but does not contain any other software which could conflict with your setup.
The latest firmware version can be downloaded from: [http://www.multitech.net/mlinux/images/mtcdt/](http://www.multitech.net/mlinux/images/mtcdt/).

The [Flashing mLinux Firmware](http://www.multitech.net/developer/software/mlinux/using-mlinux/flashing-mlinux-firmware-for-conduit/)
instructions cover both upgrading mLinux to a newer version and converting an
AEP model into a mLinux model. In both the AEP migrate and mLinux upgrade you
can use the **Using Auto-Flash During Reboot** steps. **Again, make sure to
use the `mlinux-base*.jffs2` image!**

**Important:** after flashing the device, you need to update the `opkg` cache
which can be done with the command: `opkg update`.

Example commands for upgrading / migrating:

{{<highlight bash>}}
mkdir /var/volatile/flash-upgrade
cd /var/volatile/flash-upgrade
wget -O uImage.bin http://www.multitech.net/mlinux/images/mtcdt/4.0.1/uImage--3.12.70-r20.2-mtcdt-20180620182056.bin
wget -O rootfs.jffs2 http://www.multitech.net/mlinux/images/mtcdt/4.0.1/mlinux-base-image-mtcdt-20180620182056.rootfs.jffs2
touch /var/volatile/do_flash_upgrade
reboot
{{< /highlight >}}

Then after the reboot update the `opkg` cache:

{{<highlight bash>}}
opkg update
{{< /highlight >}}

### Setting up the LoRa Gateway Bridge

1. Log in using SSH or use the USB to serial interface.

2. Download the latest `lora-gateway-bridge` `.ipk` package from:
   [https://artifacts.loraserver.io/vendor/multitech/conduit/](https://artifacts.loraserver.io/vendor/multitech/conduit/).
   Example (assuming you want to install `lora-gateway-bridge_2.5.0-r1.1_arm926ejste.ipk`):
   {{<highlight bash>}}
   wget https://artifacts.loraserver.io/vendor/multitech/conduit/lora-gateway-bridge_2.6.0-r1.1_arm926ejste.ipk
   {{< /highlight >}}

3. Now this `.ipk` package is stored on the Conduit, you can install it
   using the `opkg` package-manager utility. Example (assuming the same
   `.ipk` file):
   {{<highlight bash>}}
   opkg install lora-gateway-bridge_2.6.0-r1.1_arm926ejste.ipk
   {{< /highlight >}}

4. Update the MQTT connection details so that LoRa Gateway Bridge is able to
   connect to your MQTT broker. You will find the configuration file in the
   `/var/config/lora-gateway-bridge` directory.

5. Start LoRa Gateway Bridge and ensure it will be started on boot.
   Example:
   {{<highlight bash>}}
   /etc/init.d/lora-gateway-bridge start
   update-rc.d lora-gateway-bridge defaults
   {{< /highlight >}}

6. Be sure to add the gateway to the lora-app-server.
   See [Gateways](/lora-app-server/use/gateways/).

### Setting up the packet-forwarder

The packages installed with the commands below will by default choose the US
or EU band configuration, based on the used hardware. Please refer to the
[Multitech documentation](http://www.multitech.net/developer/software/lora/aep-lora-packet-forwarder/)
for alternative configurations.

To find out if you should follow the `MTAC-LORA-H` or `MTAC-LORA`
instructions, you could use the following commands:

{{<highlight bash>}}
# mCard in ap1 slot
mts-io-sysfs show ap1/product-id

# mCard in ap2 slot
mts-io-sysfs show ap2/product-id

# detect first mCard
mts-io-sysfs show lora/product-id
{{< /highlight >}}

#### mLinux with MTAC-LORA-H-915 or MTAC-LORA-H-868

**Important:** Follow these steps only when you have a `MTAC-LORA-H` (v1.5)
card which uses the SPI interface.

1. Log in using SSH or use the USB to serial interface.

2. Download the latest `lora-packet-forwarder` `*.ipk` package
   from [https://artifacts.loraserver.io/vendor/multitech/conduit/](https://artifacts.loraserver.io/vendor/multitech/conduit/).
   Example:
   {{<highlight bash>}}
   wget https://artifacts.loraserver.io/vendor/multitech/conduit/lora-packet-forwarder_4.0.1-r5.0_mtcdt.ipk
   {{< /highlight >}}

3. Now this `.ipk` package is stored on the Conduit, you can install it
   using the `opkg` package-manager utility. Example (assuming the same
   `.ipk` file):
   {{<highlight bash>}}
   opkg install lora-packet-forwarder_4.0.1-r5.0_mtcdt.ipk
   {{< /highlight >}}

4. As the package has the same name as the default package provided by Multitech
   you need to 'flag' the package with the status `hold` to make sure an
   `opkg upgrade` does not overwrite it:
   {{<highlight bash>}}
   opkg flag hold lora-packet-forwarder
   {{< /highlight >}}

5. Start the packet-forwarder and enable it to start on boot. Note that the
   `-ap1` or `-ap2` suffix refers to the slot in which your `MTAC-LORA-H` card
   is present. In case you have two `MTAC-LORA-H` cards, this allows you to start
   two packet-forwarder instances with each using their own configuration.
   Example:
   {{<highlight bash>}}
   /etc/init.d/lora-packet-forwarder-ap1 start
   update-rc.d lora-packet-forwarder-ap1 defaults
   {{< /highlight >}}

   **Note:** on the first start of the packet-forwarder it will detect for you
   the version of your `MTAC-LORA-H` cards (868 or 915) and if your Conduit
   has an onboard GPS. It will then automatically generate the correct
   configuration for you.

   Configuration is stored in `/var/config/lora-packet-forwarder-ap1` and
   `/var/config/lora-packet-forwarder-ap2` directories and can be modified after
   the first start.

   The build recipe of the `.ipk` package can be found at:
   [https://github.com/brocaar/loraserver-yocto](https://github.com/brocaar/loraserver-yocto).

#### mLinux with MTAC-LORA-915 or MTAC-LORA-868

**Important:** Follow these steps only when you have a `MTAC-LORA` (v1.0
card which uses the FTDI interface.

1. Log in using SSH or use the USB to serial interface.

2. Download the latest `lora-packet-forwarder-usb` `*.ipk` package
   from [https://artifacts.loraserver.io/vendor/multitech/conduit/](https://artifacts.loraserver.io/vendor/multitech/conduit/).
   Example:
   {{<highlight bash>}}
   wget https://artifacts.loraserver.io/vendor/multitech/conduit/lora-packet-forwarder-usb_1.4.1-r2.0_arm926ejste.ipk
   {{< /highlight >}}

3. Now this `.ipk` package is stored on the Conduit, you can install it
   using the `opkg` package-manager utility. Example (assuming the same
   `.ipk` file):
   {{<highlight bash>}}
   opkg install lora-packet-forwarder-usb_1.4.1-r2.0_arm926ejste.ipk
   {{< /highlight >}}

4. As the package has the same name as the default package provided by Multitech
   you need to 'flag' the package with the status `hold` to make sure an
   `opkg upgrade` does not overwrite it:
   {{<highlight bash>}}
   opkg flag hold lora-packet-forwarder-usb
   {{< /highlight >}}

5. Start the packet-forwarder and enable it to start on boot. 
   Example:
   {{<highlight bash>}}
   /etc/init.d/lora-packet-forwarder-usb start
   update-rc.d lora-packet-forwarder-usb defaults
   {{< /highlight >}}

   **Note:** on the first start of the packet-forwarder it will detect for you
   the version of your `MTAC-LORA` cards (868 or 915). It will then automatically
   generate the correct configuration for you.

   Configuration is stored in `/var/config/lora-packet-forwarder-usb` directory
   and can be modified after the first start.

   The build recipe of the `.ipk` package can be found at:
   [https://github.com/brocaar/loraserver-yocto](https://github.com/brocaar/loraserver-yocto).

#### AEP: Setting up the packet-forwarder

Use the web interface to set up the Conduit's packet forwarder.  By default, 
the connection will not be “secure” over https because the device uses a self-
signed certificate.  Accept the certificate to proceed.

1. Log in to the web-interface.
2. On the home screen, you should be able to see information about the version of the LoRa card.  Find the corresponding section on the web page:  
  http://www.multitech.net/developer/software/lora/aep-lora-packet-forwarder/  
   This page has links to basic configuration for each card version which you 
   will need below.
3. You should see the “First-Time Setup Wizard” welcome screen.  If not, you
   can access it using the menu on the left side of the screen.  Once started,
   Click “Next”.
4. Set your password on the device for the admin account and click “Next”.
5. Set the date and time and click “Next”.
6. For the “Cellular PPP Configuration,” leave all fields blank and click 
   “Next”.
7. For the “Cellular PPP Authentication”, leave the Type as NONE, and click 
   “Next”.
8. For the “IP Setup - eth0”, set up as appropriate for your network and click 
   “Next”.
9. For the “Access Configuration”, set up as appropriate for your network and 
   click “Done”.
10. On the left menu, select “Setup”, and then “LoRa Network Server” from the 
    submenu.
11. In the “LoRa Configuration” window:
12. At the top of the left column, check “Enabled”.
13. At the top of the right column, set “Mode” to be “PACKET FORWARDER”.
14. In the “Config” text box, copy and paste the configuration data for your 
    MTAC LoRa card and region.  In addition, you will want to modify/add the 
    following configuration details in the gateway\_conf section.  Leave any 
    other settings in this section as they are.  The ref\_* fields should be set 
    for the gateway.  (Altitude is specified in meters.):
    
    {{<highlight javascript>}}
    {
        ...
        "gateway_conf": {
            ...
            "server_address": "localhost",
            "serv_port_up": 1700,
            "serv_port_down": 1700,
            "fake_gps": true,
            "ref_latitude": 39.9570133,
            "ref_longitude": -105.1603241,
            "ref_altitude": 1664
        }
    }
    {{< /highlight >}}
    Note that the serv_port_up and serv_port_down represent the ports used to 
    communicate with the lora-gateway-bridge, usually on localhost (the 
    server_address parameter).  See the image above.
    
15. Select “Submit”.
16. Select the “Save and Restart” option on the left menu.

5. Be sure to add the gateway to the lora-app-server.  See [here](/lora-app-server/use/gateways/).

6. Finally, restart the system to get everything running.

### Troubleshooting

Be sure to check log files to see what is happening.  Logs can be found on the 
gateway in the directory `/var/log/`. 

Also, if the gateway seems to be running, but no statistics are 
appearing in LoRa App Server, you may be experiencing a known bug with the 
Multitech packet forwarding code.  On these systems, we need to swap out the 
application that runs for packet formwarding. The following should resolve the issue:

{{<highlight bash>}}
cd /opt/lora
mv basic_pkt_fwd-usb basic_pkt_fwd-usb.orig
ln -s gps_pkt_fwd-usb basic_pkt_fwd-usb
{{< /highlight >}}
