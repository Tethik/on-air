all: on-air

on-air: quickstart.go
	go build

install: on-air
	mkdir -p ~/on-air/
	sudo cp on-air /usr/local/bin/
	cp systemd/* ~/.config/systemd/user/
	systemctl --user daemon-reload	
	systemctl --user start onair.service
	systemctl --user enable onair.service

clean: 	
	systemctl --user stop onair.service
	systemctl --user disable onair.service
	rm ~/.config/systemd/user/onair.*
	systemctl --user daemon-reload	
	sudo rm on-air /usr/local/bin/on-air

