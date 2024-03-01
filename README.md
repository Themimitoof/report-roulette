# Report roulette

_Report roulette_ is a small binary that uses the GitLab API to elect a person in charge of writing the report of a meeting. It is also capable of automatically gathering members of a group.


## Motivation

At work, for some meetings, we need to fire a roulette to elect who's writing the meeting report. Like in all companies, people come and quit the company. So, instead of maintaining a script with a list of people, I wrote this small tool to gather the member lists of the concerned groups that are generally up-to-date.


## Installation

```
go install github.com/themimitoof/report-roulette
```

You need to generate a [Personal Access Token](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) with at least the `api` for the environment variable `GITLAB_TOKEN`.

By default, it uses the public GitLab.com instance. If you use a self-hosted one, you can specify the host by using the environment variable `GITLAB_HOST`.

```bash
export GITLAB_TOKEN=glpat-....
export GITLAB_HOST=gitlab.example.com
```

## Usage

```
→ report-roulette -help
Usage of ./report-roulette:
  -n    Don't display the list of people in the roulette.
  -s    Output the name directly instead of a phrase.

→ report-roulette @chap-api yann +pouet
Users in the roulette:
  - A... (a...)
  - A... (c...)
  - Michael Vieira (themimitoof)
  - O... (a...)
  - P... (p...)
  - pouet
  - Y... (yann)

The roulette stopped on pouet!
```

You can ask to `report-roulette` to gather the member list of one or multiple groups by using the `@` character:

```
→ report-roulette @chap-api @squad-account
Users in the roulette:
  - A... (a...)
  - A... (c...)
  - Michael Vieira (themimitoof)
  - O... (a...)
  - P... (p...)

The roulette stopped on Michael Vieira (themimitoof)!
```

You can also add/mix users:

```
→ report-roulette themimitoof gitlab_bot dummy
Users in the roulette:
  - Elmo the bot 🤖 (gitlab_bot)
  - Michael Vieira (themimitoof)
  - Y... (dummy)

The roulette stopped on D... (dummy)!
```

If for unknown reasons, an intruder joined the meeting but is not part of the company, you can add it by prepending the `+` character to its name:

```
→ report-roulette themimitoof gitlab_bot dummy +Richard
Users in the roulette:
  - Elmo the bot 🤖 (gitlab_bot)
  - Michael Vieira (themimitoof)
  - Richard
  - Y... (dummy)

The roulette stopped on Michael Vieira (themimitoof)!
```

If you want to hide the list of persons, you can use the `-n` flag and pipe the output to let a cow announce the winner:

```
→ report-roulette -n themimitoof gitlab_bot dummy +Richard | cowsay
 ______________________________________
/ The roulette stopped on Elmo the bot \
\ 🤖 (gitlab_bot)!                     /
 --------------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

You can also ask to only expose the name if you want it by using the `-s` flag:

```
→ report-roulette -n themimitoof gitlab_bot dummy +Richard
Richard
```

That's all about `report-roulette`!

## License

This project is released under the [3-Clause BSD](LICENSE) license.
