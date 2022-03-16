# GALAXIATORS MINIGAME
This code has been ported over from a private repo as an example of my work - RD

## Description
Galaxiators is an upcoming NFT game powered by the ImmutableX L2 Ethereum solution.

Whilst the main game is developed, the team wanted a minigame to provide utility to the tokens and introduce some of the background lore of the in-game universe.

In terms of codebase this was a solo project; I was the only dev and UI engineer. The storylines were constructed with the help of the writing team, and visual assets were provided by designers.

## Tech Stack
* **Frontend**: Vue3.js with the ethers lib for Web3 integration (see ./src). Also uses the vue-social-sharing lib to allow posting to social media.
* **Backend**: Golang on Google Cloud Run + Firestore (see ./backend). There is a built-in Discord notification bot to allow "Share on Discord" on the frontend.
* **Admin**: To facilitate collaborative story editing I quickly threw togther a Flask app (see ./editing)

## Status

Currently active for Galaxiator NFT holders on https://names.galaxiators.com/
