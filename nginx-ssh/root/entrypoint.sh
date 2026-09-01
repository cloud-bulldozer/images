#!/bin/sh

set -e

# Map the running UID to 'sshtest' in /etc/passwd so sshd can
# authenticate without needing setuid (OpenShift assigns random UIDs).
echo "sshtest:x:$(id -u):0:sshtest:/home/sshtest:/bin/bash" >> /etc/passwd

ssh-keygen -t rsa -f /tmp/ssh_host_rsa_key -N "" -q
ssh-keygen -t ecdsa -f /tmp/ssh_host_ecdsa_key -N "" -q
ssh-keygen -t ed25519 -f /tmp/ssh_host_ed25519_key -N "" -q

/usr/sbin/sshd -f /etc/ssh/sshd_config

exec /usr/sbin/nginx -g "daemon off;"
