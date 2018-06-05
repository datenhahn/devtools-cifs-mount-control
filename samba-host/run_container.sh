#!/bin/bash

docker run --rm --name samba-host -d -p 445:445 -t samba-host
