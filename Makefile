.PHONY: all compile docker run

all:
	@echo '                                                                       '
	@echo -e '          \e[32m,""``.\e[0m     Hello my dear user                  '
	@echo -e '         \e[32m/ _  _ \\\e[0m    Makefile for your app              '
	@echo -e '         \e[32m|(\e[0m\e[31m@\e[0m\e[32m)(\e[0m\e[31m@\e[0m\e[32m)|\e[0m    This is Octopus'
	@echo -e '         \e[32m)  \e[31m~~\e[0m  \e[32m(\e[0m    He said: use make mother faka'
	@echo -e '        \e[32m/,`))((`.\\\e[0m'
	@echo -e '       \e[32m(( ((  )) ))\e[0m'
	@echo -e '        \e[32m`\ `)(` /`\e[0m '
	@echo -e '                                                                    '
	@echo '                                                                       '
	@echo 'DEFAULT:                                                               '
	@echo '   make compile                                                        '
	@echo '   make docker                                                         '
	@echo '   '
	@echo 'RUN:'
	@echo '   make run'

compile:
#	@git config --global url."git@github.com:".insteadOf "https://github.com/"
	go get -d -u github.com/glukaza/cdn-url
	go build

docker:
	@echo 'Build Docker'
	cp cdn-url docker/rootfs/opt/
	cd docker && docker build -t ci/cdn-url:latest .

run:
	 docker run -it ci/cdn-url:latest