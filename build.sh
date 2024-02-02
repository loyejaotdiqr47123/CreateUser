#bin/sh
#设置go编译类型为windows
set GOOS=windows
#设置go编译平台为386
set GOARCH=386
# 执行go build命令
go build -o "CreateUser.exe" "CreateUser.go"
# 输出编译完成
echo "编译完成"