# Bof2PIC
BOF/COFF obj file to PIC(shellcode). by golang

```
\boftest\
测试用的bof文件和传参json文件，json文件格式与sliver一致

\loader_bin\
bofloader 核心shellcode生成

\constgen\
将核心shellcode转换成const.go文件

.\
项目主体

```

Usage:
```
.\bofgopic.exe -bof .\dir.x64.o -args .\dir.json

.\bofgopic.exe -bof .\whoami.x64.o

生成的bin文件即为PIC shellcode

```
