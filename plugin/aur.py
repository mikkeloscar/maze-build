import aur_api


def get_dependencies(pkgs, updates):
    pkgs_info = aur_api.multiinfo(pkgs)

    for pkg in pkgs_info:
        updates.add(pkg.name)

        depends = pkg.depends + pkg.make_depends
        get_dependencies(depends, updates)


updates = set()
# get_dependencies(["neovim-git"], updates)
# print(updates)
