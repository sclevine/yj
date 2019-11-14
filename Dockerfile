# pinned version of the Alpine-tagged 'go' image
FROM golang:1.13-alpine

# setup directories for the source and working
RUN mkdir -p /usr/local/src/yj/ /workdir && chown -R nobody /workdir

# copy the yj source to the image
COPY . /usr/local/src/jy/

# build and compile the jq executable
RUN cd /usr/local/src/jy && go install

# use a non-privileged user
USER nobody

# work somewhere where we can write
WORKDIR /workdir

# set the default entrypoint -- when this container is run, use this command
ENTRYPOINT [ "yj" ]

# as we specified an entrytrypoint, this is appended as an argument (i.e., `jy --help`)
CMD [ "--help" ]
