# Push `{{repo.owner}}/{{repo.name}}:{{truncate commit 7}}`
Pipeline for [branch `{{commit.branch}}` by `{{commit.author}}`]({{build.link}}): **{{build.status}}**!
> {{commit.message.title}}{{#if commit.message.body}}
>
{{{regexReplace "(?m)^" commit.message.body "> "}}}{{/if}}
