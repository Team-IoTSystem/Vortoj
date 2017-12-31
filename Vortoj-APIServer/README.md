# Vortoj-APIServer
Bonan tagon! 集められたパケットデータをAPIとして提供するモジュールです。 

## APIリファレンス

| endpoint | 説明　　　|
| -------- | -------- |
| /     | root     |
| /api     | root     |
| /api/packet/:id     |  packetテーブルのidパラメータを指定できる     |
| /api/packet/new     | packetテーブルの最新データを引き取る     |
| /api/packet/macaddress     |  packetテーブルのmacaddressを指定できる     |
| /api/distance/:id    | distanceテーブルのidパラメータを指定できる      |
| /api/distance/new    | distanceテーブルの最新データを引き取る     |
| /api/distance/macaddress   | distanceテーブルのmacaddressを指定できる,`macaddress`と`rpi_macaddress`を指定して、`new_order_one=1`にするとそれの積集合の中から最新のデータを一件引き取れる     |
