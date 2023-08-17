import sys
import ctypes
import win32net
import win32netcon

def create_user(username, password):
    user_info = {
        'name': username,
        'password': password,
        'priv': win32netcon.USER_PRIV_USER,
        'flags': win32netcon.UF_SCRIPT | win32netcon.UF_DONT_EXPIRE_PASSWD | win32netcon.UF_PASSWD_CANT_CHANGE,
    }
    
    try:
        win32net.NetUserAdd(None, 1, user_info)
    except win32net.error as e:
        print(f"Error creating user: {e}")
        return False
    return True

def add_user_to_group(username):
    try:
        win32net.NetLocalGroupAddMembers(None, 'Users', 3, [{'name': username}], 1)
    except win32net.error as e:
        print(f"Error adding user to group: {e}")
        return False
    return True

def main():
    if len(sys.argv) != 3:
        print("Usage: CreateUser.exe username password")
        return
    
    username = sys.argv[1]
    password = sys.argv[2]
    
    if create_user(username, password) and add_user_to_group(username):
        print(f"User '{username}' created and added to group 'Users'")

if __name__ == '__main__':
    main()
