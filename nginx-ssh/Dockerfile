FROM registry.access.redhat.com/ubi9/ubi-minimal
COPY root ./

RUN microdnf -y --setopt=tsflags=nodocs install nginx openssh-server openssh-clients && \
    chmod 755 /etc/nginx/ssl && \
    chmod +x /entrypoint.sh && \
    mkdir -p /etc/ssh-authorized-keys /home/sshtest && \
    chmod 755 /etc/ssh-authorized-keys /home/sshtest && \
    chmod -R g+rwX /var/lib/nginx /var/log/nginx /usr/share/nginx /run && \
    chmod g+rw /etc/passwd && \
    microdnf clean all

EXPOSE 8080 8443 2222

USER 1001

CMD ["./entrypoint.sh"]
