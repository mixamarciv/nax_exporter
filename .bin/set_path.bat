@SET this_file_path=%~dp0
@SET this_disk=%this_file_path:~0,2%

@%this_disk%
CD "%this_file_path%\.."

@SET GOROOT=d:\program\go\1.7.3\Go\
@SET GOPATH=%this_file_path%\..
@SET GIT_PATH=d:\program\git
@SET PYTHON_PATH=d:\program\Python26
@SET MINGW_PATH=c:\MINGW

@SET PATH=%GOROOT%;%GOROOT%\bin;%PATH%;
@SET PATH=%GOPATH%;%PATH%;
@SET PATH=%PYTHON_PATH%;%PATH%;
@SET PATH=%GIT_PATH%;%GIT_PATH%\bin;%PATH%;
@SET PATH=%MINGW_PATH%;%MINGW_PATH%\bin;%PATH%;
@SET PATH=%this_file_path%\..\bin;%PATH%;
