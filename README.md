# Vortoj
IoTデバイスを取り扱うあなたのソリューションです。

## What is "Vortoj"?
VortojはEsperanto語においての言の葉という意味です。
僕らの持ってるIoTデバイスを繋げ,情報を伝えて続けていくという気持ちが込められています。

## Vortoj-PacketFilter
packetをキャプチャーするモジュールです。実行して起動してください。mysqlに対してパケットを集め、プロセス間通信を行いパケットデータを後述するWebAPIサーバーに送ります

## Vortoj-APIServer
packetデータをmysqlから引き取ってWebAPIを通じて送信します。

## Vortoj-StremingAPIServer
packetデータをUnixDomainSocketのプロセス間通信を通じてデータのやり取りし,websocket通信をして提供します。

## install.sh
セットアップ用スクリプト。
12 ~ 14行目をお住まいの環境に合わせてね。

```shell
wget -O https://raw.githubusercontent.com/Team-IoTSystem/Vortoj/master/install.sh
chmod u+x install.sh
sudo -E ./install.sh
```
