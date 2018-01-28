#!/bin/bash

set -Ceu

#  Setup Vortoj-IoTSystem
#  Update: 2018/01/04
#  Raspbian Nov.2017


CMDNAME=`basename $0`

br="192.168.1.111"
ssid="Miagete-goLAN"
password="yorunohoshiwo"

while getopts e:w:h OPT
do
  case $OPT in
	b) br="$OPTARG" ;;
    p) password="$OPTARG" ;;
	s) ssid="$OPTARG" ;;      
    h) abort "Usage: $CMDNAME  [-b bridge_ip(default:$br)] [-s wifi ssid(default:$ssid) ]  [-p wifi passphrase(default:$password) ]  " ;;
  esac
done

cat <<- EOF
	your setting on
	
	bridge
		ip: $br
	
	wifis:
		ssid: $ssid
		password: $password
EOF

# Install
sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get install -y hostapd bridge-utils git vim tmux libpcap-dev mariadb-server mariadb-client python-mysqldb 
sudo systemctl stop hostapd


# Golang
wget https://storage.googleapis.com/golang/go1.9.2.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.9.2.linux-armv6l.tar.gz
sudo -u pi mkdir IoT-System
sudo chmod 765 /usr/local/go
sudo cat <<- `EOF` >> $HOME/.bashrc
	export GOROOT=/usr/local/go
	export GOPATH=$HOME/IoT-System
	export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
	`EOF`
sudo ln -s /usr/local/go/bin/go /usr/bin/go
sudo -u pi git clone --branch master --single-branch --depth=1 https://github.com/Team-IoTSystem/Vortoj.git $HOME/IoT-System
export GOROOT=/usr/local/go
export GOPATH=$HOME/IoT-System
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
go get github.com/go-sql-driver/mysql github.com/gocraft/dbr github.com/google/gopacket github.com/labstack/echo github.com/labstack/gommon github.com/gorilla/websocket github.com/dgrijalva/jwt-go


# Static IP Address
# sudo echo -e "denyinterfaces eth0 wlan0" >> /etc/dhcpcd.conf
sudo brctl addbr br0
cat <<- EOF >> /etc/dhcpcd.conf
	interface br0
	static ip_address=$br/24
	static routers=${br%.*}.1
	static domain_name_servers=${br%.*}.1
	EOF
sudo service dhcpcd reload


# Access Point
cat <<- EOF > /etc/hostapd/hostapd.conf
	interface=wlan0
	bridge=br0
	#driver=nl80211
	ssid=$ssid
	hw_mode=g
	channel=7
	wmm_enabled=0
	macaddr_acl=0
	auth_algs=1
	ignore_broadcast_ssid=0
	wpa=2
	wpa_passphrase=$password
	wpa_key_mgmt=WPA-PSK
	wpa_pairwise=TKIP
	rsn_pairwise=CCMP
	EOF
sudo sed -i -e "s@#DAEMON_CONF=\"\"@DAEMON_CONF=\"/etc/hostapd/hostapd.conf\"@" /etc/default/hostapd


# Bridge Setup
sudo brctl addif br0 eth0
cat <<- EOF >> /etc/network/interfaces
	# Bridge setup
	auto br0
	iface br0 inet manual
	bridge_ports eth0 wlan0
	EOF


#Cleanup
sudo apt-get autoremove
sudo apt-get autoclean
sudo apt-get clean
sudo rm -rf /var/lib/apt/lists/*
sudo rm -rf /var/cache/apt/archives/*
sudo reboot
