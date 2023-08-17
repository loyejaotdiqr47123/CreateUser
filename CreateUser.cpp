#include <Windows.h>
#include <Lmcons.h>
#include <iostream>

int main(int argc, char* argv[]) {
    if (argc < 3) {
        std::cout << "Usage: " << argv[0] << " <username> <password>" << std::endl;
        return 1;
    }

    LPCWSTR username = (LPCWSTR)argv[1];
    LPCWSTR password = (LPCWSTR)argv[2];

    USER_INFO_1 ui;
    DWORD dwLevel = 1;
    DWORD dwError = 0;

    ui.usri1_name = (LPWSTR)username;
    ui.usri1_password = (LPWSTR)password;
    ui.usri1_priv = USER_PRIV_USER;
    ui.usri1_home_dir = NULL;
    ui.usri1_comment = NULL;
    ui.usri1_flags = UF_SCRIPT | UF_DONT_EXPIRE_PASSWD | UF_PASSWD_CANT_CHANGE;
    
    NET_API_STATUS nStatus = NetUserAdd(NULL, dwLevel, (LPBYTE)&ui, &dwError);

    if (nStatus == NERR_Success) {
        std::cout << "User created successfully." << std::endl;
        return 0;
    } else {
        std::cerr << "Error creating user: " << nStatus << std::endl;
        return 1;
    }
}
