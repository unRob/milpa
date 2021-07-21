---
title: Why use milpa
weight: 50
---

I built milpa with a few use cases in mind:

- I need to **share groups of scripts** used by most folks in an engineering team (i.e. setup development environments, work with secrets/credentials)
- I want to share **context and code**, without having to ask folks to run stuff off READMEs
- I **forget I have all these functions and scripts** in my dotfiles/PATH/fpath that I always have to read the source for before using them, if I've not run them in a while.
- All of the above, on multiple operating systems, distribution platforms and architectures

I have used/would use milpa for:

- managing homelab services: building and deploying them, looking at their status and logs. From my personal device or CI.
- bootstrapping engineering laptops: no need to follow a README, you get the right development environment for your os/arch, your credentials setup and the code ready for you to dive in
- every-day dev workflow: connect to vpn, get credentials for a database, maybe `--connect` directly to it, abstract away all those weird cloud provider APIs and CLIs, create releases
- organize all those odd, one-off-but-not-really commands: all those shells scripts, quick ruby scripts, perl hacks and jq monstrosities can have a nice home.

milpa is great for these things, since it only requires Bash and a supported OS/arch but will setup and work with whatever your team likes using (today and beyond). milpa does the work of parsing arguments, validating them, offering help when things go sideways and so on, so these small scripts remain just so: small scripts.

## milpa is great when

- working with small programs meant to be used as either building blocks for other programs, or user-facing collections of these blocks
- milpa repos are shared over your team's source control versioning system, allowing anyone to mess with them and test them often
- writing shell scripts is the way to go. Maybe you're interacting with binaries in the system, or operating on files and directories and using your programming language of choice is overkill.

## Where milpa won't shine

I haven't tested the performance beyond dozens of repos with dozens of commands, and that being said, I can't see myself using milpa for anything more complex than that. I'd usually reach for another language to build a domain-specific CLI, specially if working with a team that is not very comfortable with shell scripting.

If there's a need to distribute stand-alone CLI programs. milpa is not the best platform to package and distribute CLIs. Complex argument parsing and auto completion are also not in my roadmap; these are very valid use cases, but keeping a simple spec that works for most cases is going to create some friction for the least-common path.


## Okay, but what about...

### Regular scripts in the filesystem

Yeah, sometimes `milpa` is overkill, specially if there's only a handful of commands that rarely change and are only used by the same few folks. `milpa` provides usability advantages over this approach which will come handy when any of the constraints listed above is not met. With `milpa`, scripts that change often get updated docs and auto completion, making onboarding easier along the way.

### Bash tools

There's some amazing tools out there, such as:

- [Bashew](https://github.com/pforret/bashew)
- [Bashinator](http://bashinator.org/)
- [Rerun](http://rerun.github.io/rerun/)
- [Sub](https://github.com/basecamp/sub)
- [Basher package manager](https://github.com/basherpm/basher)

These all serve different purposes, and some of their features overlap. milpa aims to provide the same level of support to non-bash scripts (hello applescript) without the need for another runtime to be installed. These projects are great for distributing and packaging software that may be used beyond your organization, while milpa is more focused on helping your team do your thing.

### Make

Make is great! I love using `make` for building artifacts, but I always end up hacking around it and wishing I had all of Make's strengths and the expressiveness of bash. Organizing makefiles and dealing with arguments is not something I can say I enjoy. Producing anything other than artifacts to the filesystem with Make doesn't feel right to me.

### BYOCLI

This is usually what ends up happening given a team with enough time and direction to invest in building a proper CLI using whatever language is already at use. Some teams use more than one language, which complicates this approach. It might not be easy to maintain a common interface across teams/clis, tasks like bootstrapping development environments are usually left out for the user to accomplish by reading a README, or my personal least-favorite: having to maintain another part of the system's local dev stack to run commands. For many of the tasks I use milpa with, I don't wanna reason about a programming language's abstraction over running subprocesses, where I can do fine with Bash.

