
class Node:
    def __init__(self, depends, tmp, pkg):
        self.depends = depends
        self.tmp = tmp
        self.pkg = pkg


def create_DAG(pkgs):
    """
    Create DAG from list of packages
    """
    nodes = set()
    tmp_nodes = []

def topilogical_sort(pkgs):
    L = []

