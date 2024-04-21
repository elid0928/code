import socket
import time
import os
# define IP address and port number
SERVER_HOST = '192.168.1.1'
SERVER_PORT = 8001

def receive_data(connection):
    """
    Function to receive data from the client
    """
    while True:
        data = connection.recv(1024).decode()
        if not data:
            break
        print(f"[RECEIVED] Data received from client: {data}")
    
    # close the connection
    connection.close()

def enable_server():
    """
    Function to enable server
    """
    # create a socket object
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # bind the socket to a public host, and port
    server_socket.bind((SERVER_HOST, SERVER_PORT))

    # listen for incoming connections
    server_socket.listen(1)
    print(f'[LISTENING] Server is listening on {SERVER_HOST}:{SERVER_PORT}')

    # accept connections from outside and return the connection object and client address
    connection, address = server_socket.accept()
    print(f'[NEW CONNECTION] Client {address[0]}:{address[1]} has connected.')
    os.system('/root/tcp.sh')
    # time.sleep(1000)
    # send SYN-ACK message to the client
    connection.send('SYN-ACK'.encode())
    print('[SENT] SYN-ACK message sent to client.')

    # receive ACK message from the client
    ack_message = connection.recv(1024).decode()
    if ack_message == 'ACK':
        print('[RECEIVED] ACK message received from client.')

    # receive data from the client
    receive_data(connection)

if __name__ == '__main__':
    enable_server()