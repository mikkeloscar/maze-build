# Build Archlinux packages with drone
#
#     docker build --rm=true -t mikkeloscar/maze-build .

FROM nfnty/arch-devel:latest
MAINTAINER Mikkel Oscar Lyderik <mikkeloscar@gmail.com>

# Setup build user/group
ENV UGID='1000' UGNAME='builder'
RUN \
    groupadd --gid "$UGID" "$UGNAME" && \
    useradd --create-home --uid "$UGID" --gid "$UGID" --shell /usr/bin/false "${UGNAME}"

RUN \
    # Update and install pkgbuild-introspection and jq
    pacman -Syu pkgbuild-introspection jq --noconfirm && \
    # Clean .pacnew files
    find / -name "*.pacnew" -exec rename .pacnew '' '{}' \; && \
    # Clean pkg cache
    find /var/cache/pacman/pkg -mindepth 1 -delete

# copy sudoers file
COPY contrib/etc/sudoers /etc/sudoers
# Add default mirror
COPY contrib/etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist
# Add pacman.conf template
COPY contrib/etc/pacman.conf.template /etc/pacman.conf.template
# Add gnupg config
ADD contrib/.gnupg /home/$UGNAME/.gnupg

RUN \
    chmod 'u=r,g=r,o=' /etc/sudoers && \
    chmod 'u=rwX,g=rX,o=rX' /etc/pacman.d/mirrorlist && \
    chmod 'u=rwX,g=rX,o=rX' /etc/pacman.conf.template && \
    chown "$UGNAME:$UGNAME" /etc/pacman.conf.template && \
    chown "$UGNAME:$UGNAME" -R /home/$UGNAME/.gnupg

# Add wrapper script
COPY build_wrapper.sh /usr/bin/build

# Add binary
COPY maze-build /usr/bin

ENTRYPOINT ["/usr/bin/build"]
