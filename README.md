mongodb powered blogging system in go - with an additional telnet frontend.  
uses web.go and mustache.go for http
uses mgo database driver
uses no fancy libs for telnet

(c) Leon Szpilewski
[Audio Recorder Mac](http://www.fluxforge.com/kvlt/)
Licensed under GPL v3

Dependencies:
	web.go
	mustache.go
	mgo

web/ - The http version of the blog
tnt/ - the telnet version of the blog
shared/ - shared files

how to get it running for dummies:
./getdeps.sh
./build.sh


