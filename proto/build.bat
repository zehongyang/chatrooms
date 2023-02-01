@echo off
echo Build:%1
protoc --go_out=. %1