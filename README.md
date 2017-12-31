# Vortoj
IoTデバイスを取り扱うあなたのソリューションです。

## What is "Vortoj"?
VortojはEsperanto語においての言の葉という意味です。
僕らの持ってるIoTデバイスを繋げ,情報を伝えて続けていくという気持ちが込められています。

## Vortoj-PacketFilter
packetをキャプチャーするモジュールです。実行して起動してください。mysqlに対してパケットを集め、プロセス間通信を行いパケットデータを後述するWebAPIサーバーに送ります

## Vortoj-APIServer
packetデータをmysqlから引き取ってWebAPIを通じて送信します。

## Vortoj-StremingAPIServer
packetデータをUnixDomainSocketのプロセス間通信を通じてデータのやり取りし,websocket通信をして提供します。

## prepare
セットアップ用スクリプト。
12 ~ 14行目をお住まいの環境に合わせてね。

```shell

chmod u+x prepare.sh
sudo -E ./prepare.sh
```
