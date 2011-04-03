#!/usr/bin/env python
import select
import socket
import sys
import threading
import client_thread

HOST = ''
PORT = 1337
threads = []

def serve():
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM);
    sock.bind((HOST, PORT))
    sock.listen(5)
    input = [sock, sys.stdin]
    server_run = True

    while server_run:
        inputready, outputready, exceptready = select.select(input,[],[])
        for s in inputready:
            if s == sock:
                c = client_thread.ClientThread(sock.accept())
                c.start()
                threads.append(c)
            elif s == sys.stdin:
                sin = sys.stdin.readline().strip()
                if 'quit' in sin:
                    server_run = False
    sock.close()
    for c in threads:
        c.join()
#

if __name__ == "__main__":
    serve()
