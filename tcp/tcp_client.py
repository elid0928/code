import socket

# define IP address and port number
SERVER_HOST = '34.125.33.128'
SERVER_PORT = 8001

def connect_to_server():
    """
    Function to establish a connection with the server
    """
    # create a socket object
    client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # connect to the server
    client_socket.connect((SERVER_HOST, SERVER_PORT))
    print(f'[CONNECTED] Connected to server {SERVER_HOST}:{SERVER_PORT}')

    # send SYN message to the server
    client_socket.send('SYN'.encode())
    print('[SENT] SYN message sent to server.')

    # receive SYN-ACK message from the server
    syn_ack_message = client_socket.recv(1024).decode()
    if syn_ack_message == 'SYN-ACK':
        print('[RECEIVED] SYN-ACK message received from server.')

        # send ACK message to the server
        client_socket.send('ACK'.encode())
        print('[SENT] ACK message sent to server.')
    
    # close the connection
    client_socket.close()

if __name__ == '__main__':
    connect_to_server()