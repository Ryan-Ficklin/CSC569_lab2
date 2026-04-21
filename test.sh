#!/usr/bin/env bash

# Define the commands to run
server="go run server.go"
client1="go run client.go 1"
client2="go run client.go 2"
client3="go run client.go 3"
client4="go run client.go 4"
client5="go run client.go 5"
client6="go run client.go 6"
client7="go run client.go 7"
client8="go run client.go 8"

# Opens a konsole tab with a saved layout
# Change the path to point to the layout file on your system
# KPID is used for renaming tabs

konsole --fullscreen --layout ./test_layout.json & KPID=$!

# Short sleep to let the tab creation complete
sleep 0.5

# Runs commands in Konsole panes
service="$(qdbus | grep -B1 konsole | grep -v -- -- | sort -t"." -k2 -n | tail -n 1)"

qdbus $service /Sessions/1 org.kde.konsole.Session.runCommand "${server}"
sleep 0.1
qdbus $service /Sessions/2 org.kde.konsole.Session.runCommand "${client1}"
sleep 0.1
qdbus $service /Sessions/3 org.kde.konsole.Session.runCommand "${client2}"
sleep 0.1
qdbus $service /Sessions/4 org.kde.konsole.Session.runCommand "${client3}"
sleep 0.1
qdbus $service /Sessions/5 org.kde.konsole.Session.runCommand "${client4}"
sleep 0.1
qdbus $service /Sessions/6 org.kde.konsole.Session.runCommand "${client5}"
sleep 0.1
qdbus $service /Sessions/7 org.kde.konsole.Session.runCommand "${client6}"
sleep 0.1
qdbus $service /Sessions/8 org.kde.konsole.Session.runCommand "${client7}"
sleep 0.1
qdbus $service /Sessions/9 org.kde.konsole.Session.runCommand "${client8}"

# Renames the tabs - optional
#qdbus org.kde.konsole-$KPID /Sessions/1 setTitle 1 'System Info'
#qdbus org.kde.konsole-$KPID /Sessions/2 setTitle 1 'System Monitor'

