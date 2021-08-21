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
      # {{uppercase (regexReplace "^master$" build.branch "staging")}} deployed
      **Successfully** deployed {{repo.owner}}/{{repo.name}} [`{{build.branch}}@{{truncate commit 7}}`]({{build.link}}) -> https://<URL> [[diff]({{commit.link}})]
      Author: `{{commit.author.username}}`
      > {{commit.message.title}}{{#if commit.message.body}}
      >
      {{{regexReplace "(?m)^" commit.message.body "> "}}}{{/if}}
  when:
    event: push
    status: [success]
    branch: [master, production]
```

## Developing

```sh
$ ./build.sh
$ ./test.sh
```
