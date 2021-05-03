# On Air

Connects the webcam to my "on-air" sign. Then it simply toggles our cheap wifi lamp dongle thingy (https://github.com/vikstrous/zengge-lightcontrol) on or off.

## Install

```sh
$ make
$ make install # installs the application as a systemd user service
```

## Todo

- [] Testing IRL.
- [] cli commands to turn on/off manually. And ignore for set time?

<!-- planned usage:
```sh
on-air on # turn on lamp
on-air off # turn off lamp
on-air daemon # listen for new calendar events, and turn on/off lamp
``` -->
