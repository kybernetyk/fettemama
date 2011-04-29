#!/bin/bash

mkdir deps
cd deps

#----------------------------------------------------
echo "mustache.go ..."
#git clone git://github.com/hoisie/mustache.go.git
git clone git://github.com/jsz/mustache.go.git
cd mustache.go && gomake && gomake install && cd ..
#----------------------------------------------------
echo "web.go ..."
#git clone git://github.com/hoisie/web.go.git
git clone git://github.com/jsz/web.go.git
cd web.go && gomake && gomake install && cd ..
#----------------------------------------------------
echo "mgo ..."
goinstall launchpad.net/mgo
#----------------------------------------------------
echo "done ..."

