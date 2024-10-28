FROM scratch
COPY yj /bin/yj
ENTRYPOINT ["/bin/yj"]
