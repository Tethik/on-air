# On Air

Connects my google calendar to my "on-air" sign. Uses google calendar's API (freebusy), which it polls every five minutes using a systemd timer.
Then it simply toggles our cheap wifi lamp dongle thingy (https://github.com/vikstrous/zengge-lightcontrol) on or off.

It tries to turn on the lamp around 30 seconds before the planned event. If an event is canceled, then there will be some out-of-sync as
it won't detect this until it's 5 minute timer elapses.

## Install

```
$ make
$ make install
$ make timers
```

You'll also need credentials.json, and potentially to create a token for your user.

Run `on-air` once and access the link from the console output, this will generate an access token for
the application to use on your google account.

## Usage

Just:

```
on-air
```

But I have some other plans too :)

## Todo

- [] Testing IRL.
- [] Investigate smarter usage of systemd timers (i.e. transient timers)
- [] cli commands to turn on/off manually. And ignore for set time?

<!-- planned usage:
```sh
on-air on # turn on lamp
on-air off # turn off lamp
on-air daemon # listen for new calendar events, and turn on/off lamp
``` -->
