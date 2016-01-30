Write your plugin documentation here.

The following parameters are used to configuration the plugin's behavior:

* **url** - The URL to POST the webhook to.

The following is a sample maze-build configuration in your 
.drone.yml file:

```yaml
notify:
  maze-build:
    image: mikkeloscar/maze-build
    url: http://mockbin.org/
```
