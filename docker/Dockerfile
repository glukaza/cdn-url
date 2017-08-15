FROM centos:7

MAINTAINER Alexander Donov <puzo2002@gmail.com>

RUN yum update -y && yum clean all

COPY rootfs /

EXPOSE 8080

CMD ["/opt/cdn-url"]
