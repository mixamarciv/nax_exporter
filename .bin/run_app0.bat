CALL "%~dp0/set_path.bat"


@CLS

@echo === build ===================
::go build -o app.exe

@echo ==== start ==================
app.exe cfg.json

@echo ==== end ====================
@PAUSE
