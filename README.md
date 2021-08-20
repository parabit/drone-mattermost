# About

Drone plugin for sending Mattermost notifications.

## Usage

Add as a Drone pipeline step like the following:

```yaml

steps:
- name: notify
  depends_on: [deploy]
  image: kenshaw/drone-mattermost
  settings:
    url: https://mattermost.<DOMAIN>
    token:
      from_secret: mattermost-token
    team: dev
    channel: town-square
    template: |-
      # Push `{{repo.name}}@{{build.branch}}`
      Build [{{repo.name}}]({{build.link}}) deployment {{#success build.status}}succeeded{{else}}**failed!**{{/success}}
  when:
    event: push
    status: [success, failure]
    branch: [master, production]
```

## Developing

```sh
$ ./build.sh
$ ./test.sh
```
