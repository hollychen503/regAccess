FROM alpine

ADD regAccess /regAccess
RUN chmod +x  /regAccess

EXPOSE 10080

CMD ["/regAccess" "-port=10080", "-htpasswd=auth/htpasswd"]