#!/bin/bash
#
# Run mender-client in a qemu system on your devhost
#
source .env
docker run -it \
-e SERVER_URL='https://hosted.mender.io' \
-e TENANT_TOKEN=${MENDER_TENANT_TOKEN} \
-p 8822:8822 \
mendersoftware/mender-client-qemu:latest