#!/bin/bash

# Test API methods
#   Usage:
#     bash test.sh [DAY|NIGHT|OFF]
#

HOST=hocon.dev.lan:8080

ACCEPT="Accept: application/json"
CT=""
#"Content-Type: application/json"

AUTH="Authorization: ffdsse#3"

[[ "$API_HOST" ]] || API_HOST=http://$HOST

do_cmd() {
local scene=$1
action="/api/lamp"
[[ "$scene" ]] && scene="&scene=$scene"
curl  -gs -H "$ACCEPT" "${API_HOST}$action?id=lamp1$scene"
#&scene=$scene"

#rv=$(curl -gs -H "$ACCEPT" ${API_HOST}$action?id=lamp1&scene=$scene)
#echo $rv
}

#do_cmd DAY
#sleep 1
do_cmd $1
# OFF
#sleep 1
#do_cmd NIGHT
#sleep 1
#do_cmd DAY



# | jq -r .
