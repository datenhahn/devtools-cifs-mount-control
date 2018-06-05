FROM ubuntu:xenial

RUN apt-get clean && apt-get update && apt-get install -y samba

RUN mkdir -p /sambashare && chmod 777 /sambashare

ADD smb.conf /etc/samba/smb.conf

CMD smbd -s /etc/samba/smb.conf -F
