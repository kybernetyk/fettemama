#!/bin/bash

mkdir deps
cd deps
echo "cloning mustache ..."
git clone git://github.com/hoisie/mustache.go.git
cd mustache.go && gomake && gomake install && cd ..
echo "cloning web.go ..."
git clone git://github.com/hoisie/web.go.git
cd web.go && gomake && gomake install && cd ..
echo "cloning gomongo ..."
git clone git://github.com/mikejs/gomongo.git
cd gomongo/mongo && gomake && gomake install && cd ..

