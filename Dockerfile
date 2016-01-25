# Build Archlinux packages with drone
#
#     docker build --rm=true -t mikkeloscar/drone-pkgbuild .

# FROM gliderlabs/alpine:3.2
FROM nfnty/arch-devel:latest
MAINTAINER Mikkel Oscar Lyderik <mikkeloscar@gmail.com>

# Setup build user
ENV UGID='99999' UGNAME='builder'
RUN groupadd --gid "${UGID}" "${UGNAME}"
RUN useradd --uid "${UGID}" --gid "${UGID}" --shell /usr/bin/false "${UGNAME}"

# copy sudoers file
COPY etc/sudoers /etc/sudoers
RUN chmod 'u=r,g=r,o=' /etc/sudoers

# Update and install python
RUN pacman -Syu python python-pip --noconfirm
# Clean .pacnew files
RUN find / -name "*.pacnew" -exec rename .pacnew '' '{}' \;
# Clean pkg cache
RUN find /var/cache/pacman/pkg -mindepth 1 -delete

# Add default mirror
COPY etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist
RUN chmod 'u=rwX,g=rX,o=rX' /etc/pacman.d/mirrorlist

# RUN apk-install python3
RUN mkdir -p /opt/drone
COPY requirements.txt /opt/drone/
WORKDIR /opt/drone
RUN pip install -r requirements.txt
COPY plugin /opt/drone/

USER ${UGNAME}

ENTRYPOINT ["python", "/opt/drone/plugin/main.py"]
