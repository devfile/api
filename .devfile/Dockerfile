FROM python:3-alpine
ENV GOROOT /usr/lib/go

RUN apk add --no-cache --update curl bash jq go git openssh \
&& pip3 install yq \
&& pip3 install jsonschema-cli

RUN mkdir -p /home/user/go && chmod -R a+w /home/user
ENV HOME /home/user
ENV GOPATH /home/user/go

# Set permissions on /etc/passwd and /home to allow arbitrary users to write
COPY --chown=0:0 entrypoint.sh /
RUN mkdir -p /home/user && chgrp -R 0 /home && chmod -R g=u /etc/passwd /etc/group /home && chmod +x /entrypoint.sh \
&& chgrp -R 0 /usr/lib/go && chmod -R g=u /usr/lib/go

USER 10001
WORKDIR /projects
ENTRYPOINT [ "/entrypoint.sh" ]
CMD ["tail", "-f", "/dev/null"]
