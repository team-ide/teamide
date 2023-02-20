@echo off

SET dir=%~dp0
echo dir:%dir%
cd %dir%
SET PATH=%PATH%;%dir%\lib
echo PATH:%PATH%
server.exe

pause