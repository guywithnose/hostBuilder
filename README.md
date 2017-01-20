Host Builder
A command line tool for managing frequently changing hosts files.

example hosts file:
```
127.0.0.1 www.example.com #local
10.0.0.1 www.example.com #dev
10.0.0.4 www.example.com #staging
10.0.0.5 www.example.com #awsEast
10.0.0.6 www.example.com #awsWest
```

equivalent hostBuilder config:
```
{
  "localHostnames": [],
  "ipV6Defaults": false,
  "hosts": {
    "www.example.com": {
      "current": "default",
      "options": {
        "local": "127.0.0.1",
        "dev": "10.0.0.1",
        "staging": "10.0.0.4",
        "awsEast": "10.0.0.5",
        "awsWest": "10.0.0.6"
      }
    }
  },
  "globalIPs": {},
  "groups": {}
}
```

[![asciicast](https://asciinema.org/a/7pvsjkgqy9cbdqeqo17qo6tva.png)](https://asciinema.org/a/7pvsjkgqy9cbdqeqo17qo6tva)
