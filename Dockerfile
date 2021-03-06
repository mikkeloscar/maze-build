# Build Archlinux packages with drone
#
#     docker build --rm=true -t mikkeloscar/maze-build .

FROM archlinux:latest
LABEL maintainer="Mikkel Oscar Lyderik Larsen <m@moscar.net>"

# WORKAROUND for glibc 2.33 and old Docker
# See https://github.com/actions/virtual-environments/issues/2658
# Thanks to https://github.com/lxqt/lxqt-panel/pull/1562
RUN curl -fsSL \
    "https://repo.archlinuxcn.org/x86_64/glibc-linux4-2.33-4-x86_64.pkg.tar.zst" \
    | bsdtar -C / -xvf -

RUN \
    # Update and install packages
    pacman -Syu \
        base-devel \
        git \
        jq \
        --noconfirm && \
    # Clean .pacnew files
    find / -name "*.pacnew" -exec rename .pacnew '' '{}' \;

# Setup build user/group
ENV UGID='1000' UGNAME='builder'
RUN \
    groupadd --gid "$UGID" "$UGNAME" && \
    useradd --create-home --uid "$UGID" --gid "$UGID" --shell /usr/bin/false "${UGNAME}"

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

# set git environment
ENV GIT_AUTHOR_EMAIL=maze-build GIT_AUTHOR_NAME=maze-build \
    GIT_COMMITTER_EMAIL=maze-build GIT_COMMITTER_NAME=maze-build

# Add binary
COPY build/linux/maze-build /usr/bin

ENTRYPOINT ["/usr/bin/maze-build"]
