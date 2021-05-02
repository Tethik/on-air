# On Air

Connects my google calendar to my "on-air" sign. Uses google calendar's API (freebusy), and polls every five minutes.

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

<!-- planned usage:
```sh
on-air on # turn on lamp
on-air off # turn off lamp
on-air daemon # listen for new calendar events, and turn on/off lamp
``` -->
