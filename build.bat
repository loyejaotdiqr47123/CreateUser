@echo off
cd build && "C:\Program Files (x86)\Windows Kits\10\bin\10.0.19041.0\x64\cl.exe" /LD CreateUserDLL.cpp /link /OUT:CreateUserDLL.dll netapi32.lib
