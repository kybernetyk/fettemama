#!/usr/bin/env python
''' 
    client_thread.py:
        threading handler for client session
    (c) Leon Szpilewski 2011
        http://nntp.pl
    License: GPL v3
'''

import socket
import sys
import threading
import blog

#oh wee, python threading suxx. need to build a class >.<
class ClientThread(threading.Thread):
    def __init__(self, (client, address)):
        threading.Thread.__init__(self)
        self.client = client
        self.address = address
        self.size = 1024
        self.blog = blog.Blog()

    def run(self):
        s = self.blog.render_version()
        s += 'help for help\n'
        s += self.blog.render_prompt()
        self.client.send(s)
        while 1:
            data = self.client.recv(self.size).strip("\r\n")
            if data:
                #self.client.send("LOL THX UR: " + data)
                command, ret = self.blog.process_input(data)
                ret += self.blog.render_prompt()
                self.client.send(ret)

                if command == 'close':
                    self.client.close()
                    break
            else:
                self.client.close();
                break;
#


