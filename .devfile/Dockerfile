#
#
# Copyright Red Hat
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
FROM python:3-alpine
ENV GOROOT /usr/lib/go

RUN apk add --no-cache --update curl bash jq go git openssh \
&& pip3 install yq \
&& pip3 install pyyaml \
&& pip3 install jsonschema

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
