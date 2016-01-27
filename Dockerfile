# Build Archlinux packages with drone
#
#     docker build --rm=true -t mikkeloscar/drone-pkgbuild .

FROM nfnty/arch-devel:latest
MAINTAINER Mikkel Oscar Lyderik <mikkeloscar@gmail.com>

# Setup build user/group
ENV UGID='99999' UGNAME='builder'
RUN groupadd --gid "$UGID" "$UGNAME"
RUN useradd --uid "$UGID" --gid "$UGID" --shell /usr/bin/false "${UGNAME}"

# copy sudoers file
COPY etc/sudoers /etc/sudoers
RUN chmod 'u=r,g=r,o=' /etc/sudoers

# Update and install pkgbuild-introspection
RUN pacman -Syu pkgbuild-introspection --noconfirm
# Clean .pacnew files
RUN find / -name "*.pacnew" -exec rename .pacnew '' '{}' \;
# Clean pkg cache
RUN find /var/cache/pacman/pkg -mindepth 1 -delete

# Add default mirror
COPY etc/pacman.d/mirrorlist /etc/pacman.d/mirrorlist
RUN chmod 'u=rwX,g=rX,o=rX' /etc/pacman.d/mirrorlist

# Add pacman.conf template
COPY etc/pacman.conf.template /etc/pacman.conf.template
RUN chmod 'u=rwX,g=rX,o=rX' /etc/pacman.conf.template
RUN chown "$UGNAME:$UGNAME" /etc/pacman.conf.template

# Add binary
COPY drone-pkgbuild /usr/bin

USER $UGNAME

ENTRYPOINT ["/usr/bin/drone-pkgbuild"]
