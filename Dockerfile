FROM scratch

COPY ./alertmanager-telegram /
ENTRYPOINT ["/alertmanager-telegram"]
CMD [ "daemon" ]
