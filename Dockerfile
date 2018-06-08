FROM alpine

ADD regAccess /usr/bin/
RUN chmod +x  /usr/bin/regAccess

EXPOSE 10080

WORKDIR /

#CMD ["regAccess" "-port=10080", "-htpasswd=auth/htpasswd"]

ENTRYPOINT ["/usr/bin/regAccess", "-port=10080", "-htpasswd=auth/htpasswd" ]