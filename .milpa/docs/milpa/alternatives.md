---
title: Why use milpa
description: Use-cases and alternatives when using milpa
weight: 2
---

I built `milpa` with a few use cases in mind:

- To **share groups of scripts** used by a group of folks in an engineering team (i.e. setup development environments, work with secrets/credentials)
- to **share context and code**, without having to ask folks to run stuff off READMEs (i.e. pull service logs, forward ports), and
- to **quickly build scripts** that will be documented sufficiently well for my forgetful future-self.

I have used/would use `milpa` for:

- managing homelab services: building and deploying them, looking at their status and logs. From my personal device or CI.
- bootstrapping engineering laptops: no need to follow a README, you get the right development environment for your os/arch, your credentials setup and the code ready for you to dive in.
- every-day dev workflow: connect to vpn, get credentials for a database, maybe `--connect` directly to it, abstract away all those weird cloud provider APIs and CLIs, create releases.
- organize all those odd, one-off-but-not-really commands: all those shells scripts, quick ruby scripts, perl hacks and jq monstrosities can have a nice home.

`milpa` is great for these things, since it only requires Bash and a supported OS/arch but will setup and work with whatever your team likes using (today and beyond). `milpa` does the work of parsing arguments, validating them, offering help when things go sideways and so on, so these small scripts remain just so: small scripts.

See [unRob/nidito](https://github.com/unRob/nidito/tree/master/.milpa) for a working example, and/or the drop that that spilled the cup and got me to work on `milpa`.

## `milpa` is great for

- working with small programs meant to be used as either building blocks for other programs, or user-facing collections of these blocks
- sharing related commands using your team's source control versioning system, allowing anyone to mess with them and test them often
- when writing shell scripts is the way to go. Maybe you're interacting with binaries in the system, or operating on files and directories and using your programming language of choice is overkill.

Tasks like bootstrapping development environments are usually left out for the user to accomplish by following through a README or wiki with likely outdated links; `milpa` is great when these tasks can be accomplished by prompting for information, querying identity providers and then running configuration commands or modifying the filesystem directly.

My aim is to make following the [Command Line Interface Guidelines](https://clig.dev) require as little effort as possible.

## Where `milpa` won't shine

I haven't tested the performance beyond dozens of repos with dozens of commands, and that being said, I can't see myself using `milpa` for anything more complex than that. I'd usually reach for another language to build a domain-specific CLI, specially if working with a team that is not very comfortable with shell scripting.

When there's a need to distribute stand-alone CLI programs. `milpa` is not the best platform to package and distribute CLIs. Complex argument parsing and auto completion are also not in my roadmap; these are very valid use cases, but keeping a simple spec that works for most cases is going to create some friction for the least-common path.

When the primary runner of commands is not gonna be a human. Sorry robot friends! `milpa`'s features are oriented primarily towards improving the experience of maintaining and running scripts by humans, and while there's nothing wrong with having an automated system (say CI, for example) run `milpa` commands, there is an overhead to consider by invoking `milpa` (both cognitive overhead and in terms of resource usage).

## Okay, so what about...

### Regular scripts in the filesystem

These are fine, as long as you know the path to your script, it doesn't consume many arguments/options, it rarely changes and/or is used only by the same few folks. `milpa` provides usability advantages over this approach which will come handy when any of the constraints listed before are not met. With `milpa`, scripts that change often get updated docs and auto completion for free, and may be appreciated by new users and non-regulars alike.

### Bash tools

There's some amazing tools out there, such as:

- [Bashew](https://github.com/pforret/bashew)
- [Bashinator](http://bashinator.org/)
- [Rerun](http://rerun.github.io/rerun/)
- [Sub](https://github.com/basecamp/sub)
- [Basher package manager](https://github.com/basherpm/basher)

These all serve different purposes, and there is some feature overlap. `milpa` aims to provide the same level of support to non-bash scripts (hello applescript) without the need for another runtime to be installed. These projects are great for distributing and packaging software that may be used beyond your organization, while `milpa` aims to help smaller user bases, such as engineering teams, keep their shared scripts organized and up-to-date.

### Make

Make is great! I love using `make` for building artifacts, but I always end up hacking around it and wishing I had all of Make's strengths and the expressiveness of bash. I used `make` to build `milpa` initially since it's really a great tool, has a huge user base and flavors of it produce software used by millions. That being said, organizing makefiles and dealing with arguments is not something I can say I enjoy, and producing anything other than artifacts to the filesystem with Make doesn't feel right to me.

### BYOCLI

This is usually what ends up happening given a team with enough time and direction to invest in building a proper CLI using whatever language is already at use. Some teams use more than one language, which may complicate this approach. In my limited experience with engineering teams of less than 300 folks, many of the bootstrapping tasks will involve operations that can easily be expressed with a shell scripting language and often involve less code compared to more expressive languages at use for building applications.
