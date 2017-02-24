# gits
#### a simple GitHub hook deployment application

This project was made out of the need for simplicity. Moving away from a very
clunky Jenkins install to something much lighter. All it does is listen for
GitHub hooks, and runs some given commands.

Right now you configure GitHub hooks manually (with a secret) and gits will
listen for any hook related to a configured project.

### Installation
Run the [latest binary file](https://github.com/btkostner/gits/releases).

### Configuration
gits expects a configuration file in `yaml` format. Here is an example one:
```
---

log:
  level: debug

server:
  port: 4200

projects:
  btkostner/gits:
    secret: secretsecret
    commands:
      push:
        - echo "{{ branch }}"
```


### License
MIT
