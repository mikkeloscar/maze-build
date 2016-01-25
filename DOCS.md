Write your plugin documentation here.

The following parameters are used to configuration the plugin's behavior:

* **url** - The URL to POST the webhook to.

The following is a sample drone-pkgbuild configuration in your 
.drone.yml file:

```yaml
notify:
  drone-pkgbuild:
    image: mikkeloscar/drone-pkgbuild
    url: http://mockbin.org/
```
