#!/bin/bash
echo "building telnet ..."
cd tnt && gomake clean && gomake && cd ..
echo "building web ..."
cd web && gomake clean && gomake && cd ..
echo "done"
