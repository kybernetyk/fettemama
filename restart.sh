#!/bin/bash

# warning: this restart script is custom to my server 
# (the blog runs behind an apache that needs to be restarted too)

echo "telnet blog"
killall telblog
sleep 2
cd tnt 
./telblog > /dev/null &
cd ..

echo "web blog"
killall blog
sleep 2
apache2ctl restart
cd web
./blog > /dev/null &
cd ..

echo "done"
