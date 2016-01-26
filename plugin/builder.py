import os
from srcinfo.parse import parse_srcinfo

class Builder:
    def __init__(self, workdir):
        self.workdir = workdir

    def build(self, pkgs):
        self._update()
        # TODO check version on git packages
        # TODO: topological sort.

        archives = []
        for pkg in pkgs:
            archives += self._build_pkg(pkg)
        return archives

    def _update_pkg_src(self, pkg_path):
        subprocess.Popen(["makepkg", "-os", "--noconfirm"], cwd=pkg_path)
        subprocess.Popen(["mksrcinfo"], cwd=pkg_path)

        file_path = os.path.join(pkg_path, ".SRCINFO")
        return get_srcinfo(file_path)


    def _update(self):
        subprocess.run(["sudo", "pacman", "-Syu", "--noconfirm"])

    def _build_pkg(self, pkg):
        path = os.path.join(self.workdir, pkg)
        subprocess.Popen(["makepkg", "-is", "--noconfirm"], cwd=path)

        pkgs = []
        for f in os.listdir(path):
            if f.endswidth("pkg.tar.xz")
                pkgs.append(os.path.join(path, f))

        return pkgs


def get_srcinfo(path):
    with open(file_path, "r") as f:
        srcinfo = f.read()
        (info, errors) = parse_srcinfo(srcinfo)
        if errors:
            # TODO raise error
            pass
        return info
