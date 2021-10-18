 windres -o main-res.syso main.rc
 go build -o bin/wshub.exe
 "C:\Tools\signtool.exe" sign /f "E:\Github\ReopeCertificate.cer" /tr http://timestamp.comodoca.com /td sha256 bin/wshub.exe