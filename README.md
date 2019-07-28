# `humble-library` WIP
[![CircleCI](https://circleci.com/gh/wizardsoftheweb/humble-library/tree/dev.svg?style=svg)](https://circleci.com/gh/wizardsoftheweb/humble-library/tree/dev) [![Coverage Status](https://coveralls.io/repos/github/wizardsoftheweb/humble-library-go/badge.svg?branch=dev)](https://coveralls.io/github/wizardsoftheweb/humble-library-go?branch=dev)

The repo right now is basically a massive PoC and probably isn't useful until it becomes a CLI.

I've tidied much of the PoC up but I haven't enabled downloads yet.

## Overview

Humble has a lot of trouble with larger libraries on all platforms. Sidestepping its constraints is fairly easy thanks to [Hayden Schiff](https://www.schiff.io/) and [his great writeup on reverse-engineering the mobile app](https://www.schiff.io/blog/2017/07/21/reverse-engineering-humble-bundle-api). In the two years since that was published, there have been a few API changes that have to be taken into account, but nothing too hard to figure out.
