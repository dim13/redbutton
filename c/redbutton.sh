#!/bin/sh

if [ $# -ne 1 ]; then
	exit 1
fi

case $1 in
	armed)
		echo armed
		;;
	reset)
		echo reset
		;;
	launch)
		echo fire
		;;
	locked)
		echo lock
		;;
	*)
		echo $1
esac
