#include <Windows.h>
#include <Lmcons.h>

extern "C" __declspec(dllexport) DWORD CreateUser(LPCWSTR username, LPCWSTR password) {
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
    return nStatus;
}
