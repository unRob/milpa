---
title: Why use milpa?
description: Use-cases and alternatives when using milpa
weight: 2
---

I built `milpa` with a few use cases in mind:

- To **share groups of scripts** like those needed by folks in an engineering team (i.e. setup development environments, work with secrets/credentials)
- to **share context and code**, without having to ask folks to run stuff off READMEs (i.e. pull service logs, forward ports), and
- to **quickly build scripts** that will be documented sufficiently well for my forgetful future-self.

My goal with `milpa` is to make following the [Command Line Interface Guidelines](https://clig.dev) require as little effort as possible.


## `milpa` is great for

- working with small programs meant to be used as either building blocks for other programs, or user-facing collections of these blocks
- sharing related commands using your team's source control versioning system, allowing anyone to mess with them and test them often
- when writing shell scripts is the way to go. Maybe you're interacting with binaries in the system, or operating on files and directories and using your programming language of choice is overkill.
- turning runbooks into actionable, runnable scripts. Up-to-date remediation one `git pull && milpa runbook ...` away.

Tasks like bootstrapping development environments are usually left out for the user to accomplish by following through a README or wiki with likely outdated links; `milpa` is great when these tasks can be accomplished by prompting for information, querying identity providers and then running configuration commands or modifying the filesystem directly.

> ⚠️ `milpa` is still under development and is currently on alpha!

### Examples

Here's some examples of how I've used used `milpa` so far:

- for **managing homelab services** (i.e. [unRob/nidito](https://github.com/unRob/nidito/tree/master/.milpa)): building and deploying them, looking at their status and logs. From my personal device or CI.
- **bootstrapping engineering laptops** (i.e. [unRob/dotfiles](https://github.com/unRob/dotfiles/tree/master/.milpa/commands/computar)): no need to follow a README, you get the right development environment for your os/arch, your credentials setup and the code ready for you to dive in.
- **every-day dev workflow** (i.e. [unRob/milpa](https://github.com/unRob/milpa/tree/main/repos/internal/commands/)): lint and test a codebase, connect to vpn, get credentials to resources, maybe `--connect` to them, abstract away APIs (internal, cloud provider and SaaS) and CLIs, toggle feature gates, build reports and update google sheets with the results.
- as a way to organize all those odd, one-off-but-not-really commands: found a nice home for [shell scripts](https://github.com/unRob/dotfiles/blob/master/.milpa/commands/code/todo.sh), quick ruby scripts, perl hacks and [jq monstrosities](https://github.com/unRob/dotfiles/blob/master/.milpa/commands/creds.sh) that used to live in my `~/.zsh_history`.

I found `milpa` useful for these particular problems, since it only requires Bash and a supported OS/arch. `milpa` helps setup the computers I use (directly or remotely) and makes working with whatever languages my employer prefers (today and beyond) a breeze to integrate into a cohesive set of scripts. `milpa` takes care of parsing, validation, offering help when things go sideways, and so on, so these small scripts remain just so: _small scripts_.

## Where `milpa` won't shine

I haven't tested the performance beyond dozens of repos with dozens of commands, and that being said, I can't see myself using `milpa` for anything more complex than that. I'd usually reach for another language to build a domain-specific CLI, specially if working with a team that is not very comfortable with shell scripting.

When there's a need to distribute stand-alone CLI programs, `milpa` won't be the best method to package and distribute CLIs. While facilities to work with `milpa` repos exists (see [`milpa itself repo install`](/.milpa/commands/itself/repo/install)), it may be less than ideal since there's a dependency in `milpa`.

`milpa` could be the wrong tool when the primary runner of commands is not gonna be a human. Sorry robot friends! `milpa`'s features are oriented primarily towards improving the experience of maintaining and running scripts by humans, and while there's nothing wrong with having an automated system (say CI, for example) run `milpa` commands, there is an overhead to consider by invoking `milpa` (both cognitive overhead and in terms of resource usage).

That being said, it can be done and it works fine (or, if you like writing bash scripts as much as I do, _beautifully_); check out [unRob/smoked-by-the-house](https://github.com/unRob/smoked-by-the-house) where `milpa` orchestrated and operated a Raspberry Pi running an art installation for three months of 2022 at the Anahuacalli Museum in Mexico City.


## Alternatives to `milpa`

### Regular scripts in the filesystem

These are great, as long as: a) you know the path to your script, b) it doesn't consume many arguments/options, c) rarely changes, and/or d) is used only by the same few folks.

`milpa` provides usability advantages over this approach which will come handy when any of the constraints listed before are not met. With `milpa`, scripts that change often get updated docs and autocomplete for free, and may be appreciated by new users and non-regulars alike. Coming back to rarely used commands is one `--help` away when these scripts are part of a `milpa` repo.

### Bash tools

There's some amazing tools out there, such as:

- [Bashew](https://github.com/pforret/bashew)
- [Bashinator](http://bashinator.org/)
- [Rerun](http://rerun.github.io/rerun/)
- [Sub](https://github.com/basecamp/sub)
- [Basher package manager](https://github.com/basherpm/basher)
- [Criteo's command-launcher](https://github.com/criteo/command-launcher)

These all serve different purposes, and there is some feature overlap between `milpa` and each of these.

`milpa` aims to provide the same level of support to non-bash scripts (hello applescript) without the need for another runtime to be installed. These projects are great for distributing and packaging software that may be used beyond your organization, while `milpa` aims to help smaller user bases, such as engineering teams, keep their shared scripts organized and up-to-date.

### Make

Make is great! I love using `make` for building artifacts, but given enough time, I'll always end up hacking around it and wishing I had all of Make's strengths and the expressiveness of bash. I used `make` to build `milpa` initially since it's really a great tool, has a huge user base and flavors of it produce software used by millions.

That being said, organizing makefiles and dealing with arguments is not something I can say I enjoy, and producing anything other than artifacts to the filesystem with Make doesn't feel right to me. `milpa` can be a good companion to Make workflows, and fill in the gaps upstream or downstream of `make`.

### BYOCLI

Building your own CLI is usually what ends up happening given a team with enough time and direction to invest in building a proper CLI with whatever language is already at use. Some teams use more than one language, which may complicate this approach. In the microservices world, many codebases come with their own CLIs that may follow slightly different conventions, are seldomly documented and often just end up calling other binaries through `exec`.

In my limited experience with engineering teams of less than 300 folks, many of the bootstrapping tasks will involve operations that can easily (and more succinctly) be expressed with a shell scripting language. `milpa` could also be a useful intermediate step, that could help teams avoid the dread of jira runbooks until it's a good time to build your own CLI.

Building a CLI is something that happens to me somewhat often, and when golang is a good choice, I get most of `milpa`'s niceties by building my CLI on [chinampa](https://git.rob.mx/nidito/chinampa). `joao`, a configuration manager, is an example of something that started as a [milpa repo](https://github.com/unRob/nidito/tree/0812e0caf6d81dd06b740701c3e95a2aeabd86de/.milpa/commands/nidito/config) that later became [it's own CLI](https://git.rob.mx/nidito/joao); getting it to work with `milpa` first was fundamental in figuring out _what_ it needed to do, and it only took a few hours.
