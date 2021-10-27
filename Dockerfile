FROM ubuntu:20.04
MAINTAINER Sean Williams <sean@xeed.io>

#Setup build environment for libpam
RUN apt-get update -y
RUN apt-get upgrade -y
RUN apt-get -y install \
      libpam-modules \
      libpam-modules-bin \
      pamtester

COPY https-user-management.deb /tmp/

RUN dpkg -i /tmp/https-user-management.deb && \
      rm -f /tmp/https-user-management.deb
