# Build Archlinux packages with drone
#
#     docker build --rm=true -t mikkeloscar/drone-pkgbuild .

FROM nfnty/arch-devel:latest
MAINTAINER Mikkel Oscar Lyderik <mikkeloscar@gmail.com>

# Setup build user/group
ENV UGID='1000' UGNAME='builder'
RUN \
    groupadd --gid "$UGID" "$UGNAME" && \
    useradd --uid "$UGID" --gid "$UGID" --shell /usr/bin/false "${UGNAME}"

RUN \
    # Update and install pkgbuild-introspection
    pacman -Syu pkgbuild-introspection --noconfirm && \
    # Clean .pacnew files
    find / -name "*.pacnew" -exec rename .pacnew '' '{}' \; && \
    # Clean pkg cache
    find /var/cache/pacman/pkg -mindepth 1 -delete

# copy sudoers file
COPY etc/sudoers /etc/sudoers
# Add default mirror
COPY etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist
# Add pacman.conf template
COPY etc/pacman.conf.template /etc/pacman.conf.template

RUN \
    chmod 'u=r,g=r,o=' /etc/sudoers && \
    chmod 'u=rwX,g=rX,o=rX' /etc/pacman.d/mirrorlist && \
    chmod 'u=rwX,g=rX,o=rX' /etc/pacman.conf.template && \
    chown "$UGNAME:$UGNAME" /etc/pacman.conf.template

# Add binary
COPY drone-pkgbuild /usr/bin

USER $UGNAME

ENTRYPOINT ["/usr/bin/drone-pkgbuild"]
