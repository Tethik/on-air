all: on-air

on-air: cmd/*.go lib/*.go
	go build

install: on-air
	go install
	cp default-config.yaml ~/.on-air.yaml

install-service:
	cp systemd/* ~/.config/systemd/user/
	systemctl --user daemon-reload	
	systemctl --user start on-air.service
	systemctl --user enable on-air.service

clean: 	
	systemctl --user stop on-air.service &>/dev/null
	systemctl --user disable on-air.service &>/dev/null
	rm ~/.config/systemd/user/on-air.* &>/dev/null
	systemctl --user daemon-reload	&>/dev/null
	go clean -i
	

