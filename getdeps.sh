#!/bin/bash

mkdir deps
cd deps
echo "mustache.go ..."
git clone git://github.com/hoisie/mustache.go.git
cd mustache.go && gomake && gomake install && cd ..
echo "web.go ..."
git clone git://github.com/hoisie/web.go.git
cd web.go && gomake && gomake install && cd ..
echo "gomongo ..."
goinstall github.com/mikejs/gomongo/mongo

