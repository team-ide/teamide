@echo off

cd ../../teamide-html

:: set NODE_OPTIONS=--openssl-legacy-provider%%
set NODE_OPTIONS=--openssl-legacy-provider
npm run serve

pause
