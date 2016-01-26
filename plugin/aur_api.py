import requests

AUR_URL = "https://aur.archlinux.org/rpc.php?"

def get(payload):
    payload["v"] = 4
    r = requests.get(AUR_URL, params=payload)
    json = r.json()
    if "error" in json:
        return None
    return json

def search(query):
    payload = { "type": "search", "arg": query }
    json = get(payload)
    pkgs = []
    if json:
        for result in json["results"]:
            pkgs.append(Package(result))
    return pkgs

def msearch(user):
    payload = { "type": "msearch", "arg": user }
    json = get(payload)
    pkgs = []
    if json:
        for result in json["results"]:
            pkgs.append(Package(result))
    return pkgs

def info(pkg):
    payload = { "type": "info", "arg": pkg }
    json = get(payload)
    if json:
        return Package(json["results"][0])
    return None

def multiinfo(pkgs):
    payload = { "type": "multiinfo", "arg[]": pkgs }
    json = get(payload)
    pkgs = []
    if json:
        for result in json["results"]:
            pkgs.append(Package(result))
    return pkgs

def set_value(stack, key, default):
    return default if key not in stack else stack[key]

class Package:
    def __init__(self, json):
        self.id = set_value(json, "ID", 0)
        self.name = set_value(json, "Name", "")
        self.package_base_id = set_value(json, "PackageBaseID", 0)
        self.package_base = set_value(json, "PackageBase", "")
        self.version = set_value(json, "Version", "")
        self.description = set_value(json, "Description", "")
        self.version = set_value(json, "URL", "")
        self.num_votes = set_value(json, "NumVotes", 0)
        self.popularity = set_value(json, "Popularity", 0.0)
        self.out_of_date = set_value(json, "OutOfDate", None)
        self.maintainer = set_value(json, "Maintainer", "")
        self.first_submitted = set_value(json, "FirstSubmitted", 0)
        self.last_modified = set_value(json, "LastModified", 0)
        self.url_path = set_value(json, "URLPath", "")
        self.license = set_value(json, "License", "")
        self.depends = set_value(json, "Depends", [])
        self.make_depends = set_value(json, "MakeDepends", [])
        self.otp_depends = set_value(json, "OptDepends", [])
        self.provides = set_value(json, "Provides", [])
        self.conflicts = set_value(json, "Conflicts", [])

    def __repr__(self):
        return "<%s: %s>" % (type(self).__name__, self.name)
