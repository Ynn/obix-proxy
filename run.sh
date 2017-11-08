#!/bin/sh
echo "STARTING PROXY FOR $OBIX_SERVER_URL"
(cd /app && go build proxy.go)
/app/proxy
exit