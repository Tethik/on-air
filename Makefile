all: on-air

on-air: quickstart.go
	go build

install: on-air
	mkdir ~/on-air/
	sudo cp on-air /usr/local/bin/

clean: 
	sudo rm on-air /usr/local/bin/on-air
	systemctl stop onair.timer
	systemctl disable onair.timer
	rm ~/.config/systemd/user/onair.*
	systemctl daemon-reload	

timers:
	cp systemd/* ~/.config/systemd/user/
	systemctl --user daemon-reload	
	systemctl --user start onair.timer
	systemctl --user enable onair.timer