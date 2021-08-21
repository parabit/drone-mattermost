# {{uppercase (regexReplace "^master$" build.branch "staging")}} deployed
**Successfully** deployed {{repo.owner}}/{{repo.name}} [`{{build.branch}}@{{truncate commit 7}}`]({{build.link}}) -> https://<URL> [[diff]({{commit.link}})]
Author: `{{commit.author.username}}`
> {{commit.message.title}}{{#if commit.message.body}}
> 
{{{regexReplace "(?m)^" commit.message.body "> "}}}{{/if}}
