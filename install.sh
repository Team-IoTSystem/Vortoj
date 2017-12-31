#!/bin/bash

set -Ceu

#  Setup Vortoj-IoTSystem
#  Update: 2017/12/23
#  Raspbian Nov.2017


CMDNAME=`basename $0`

br="192.168.10.111"
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
#TODO:洗い出し
sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get install -y hostapd bridge-utils git
sudo systemctl stop hostapd
sudo debconf-set-selections <<< 'mysql-server-5.5 mysql-server/root_password password P@ssw0rd'
sudo debconf-set-selections <<< 'mysql-server-5.5 mysql-server/root_password_again password P@ssw0rd'
sudo apt-get -y install mysql-server

cat <<- EOF >> /etc/my.conf
	[mysqld]
	character-set-server=utf8
	[client]
	default-character-set=utf8
	EOF


# Go環境
#(多分色々足りてない)
wget https://storage.googleapis.com/golang/go1.9.2.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.9.2.linux-armv6l.tar.gz
mkdir IoT-System
sudo chmod 765 /usr/local/go
sudo cat <<- `EOF` >> $HOME/.bashrc
	export GOROOT=/usr/local/go
	export GOPATH=$HOME/IoT-System
	export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
	`EOF`


sudo brctl addbr br0
# Static IP Address
# echo -e "net.ifname=0" | sudo tee /boot/cmdline.txt
# sudo echo -e "denyinterfaces eth0 wlan0" >> /etc/dhcpcd.conf
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
sudo sed -i -e "s@#DAEMON_CONF=""@DAEMON_CONF=\"/etc/hostapd/hostapd.conf\"@" /etc/default/hostapd


# Bridge Connection
sudo brctl addif br0 eth0
cat <<- EOF >> /etc/network/interfaces
	# Bridge setup
	auto br0
	iface br0 inet manual
	bridge_ports eth0 wlan0
	EOF


#TODO:githubから $HOME/IoT-System へ持ってくる処理↓


#Cleanup
sudo apt-get autoremove
sudo apt-get autoclean
sudo apt-get clean
sudo rm -rf /var/lib/apt/lists/*
sudo rm -rf /var/cache/apt/archives/*
sudo reboot
