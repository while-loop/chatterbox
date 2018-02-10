FROM alpine
COPY chatterbox /usr/local/bin
RUN chmod +x /usr/local/bin/chatterbox
CMD ["/usr/local/bin/chatterbox"]