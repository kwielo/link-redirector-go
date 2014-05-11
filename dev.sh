#!/bin/bash

pwd=$(pwd)
rootpath=/usr/local/go
path=$pwd

GOsetenv() {
	echo "Setting GOPATH to: $path"
	export GOROOT=$rootpath
	export GOPATH=$path
}

GOinstall() {
	if [ "$GOPATH" != "$path" ]; then
		GOsetenv
	fi
	echo "Intalling go script"
	cd src/fakje/gogogo
	go install
}

GOrun() {
	GOinstall
	echo "Running script"
	echo "---"
	cd $pwd/bin/
	./gogogo
}

case $1 in
	install)
		GOinstall
	;;
	run)
		GOrun
	;;
	setup)
		GOsetenv
	;;
	*)
		echo "Commands:"
		echo "dev [ install | run | setup ]"
	;;
esac


