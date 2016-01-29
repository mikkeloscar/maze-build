Write your plugin documentation here.

The following parameters are used to configuration the plugin's behavior:

* **url** - The URL to POST the webhook to.

The following is a sample maze-drone configuration in your 
.drone.yml file:

```yaml
notify:
  maze-drone:
    image: mikkeloscar/maze-drone
    url: http://mockbin.org/
```
