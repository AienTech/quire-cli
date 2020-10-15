# Quire CLI

This is a CLI application made for [Quire.io](https://quire.io), to make it easier to work with different features of the Quire App.

## Features
| Description                           	| Status 	|
|---------------------------------------	|--------	|
| Authorization                         	| Stable 	|
| Create a GIT branch based on a ticket 	| Stable 	|
| Finish a GIT branch for a ticket      	| Stable    |
| [Other features will come soon]       	| -      	|

## How to install

Simply run the following command
```
curl -s https://raw.githubusercontent.com/AienTech/quire-cli/master/install.sh | sh 
```

### To uninstall
also run
```
curl -s https://raw.githubusercontent.com/AienTech/quire-cli/master/uninstall.sh | sh 
```

## How to use?
The first step you have to go through is to authorize the app to be connected to you quire account. For that run the following command
```
$ quire authorize
```
This will open a browser window and lets you authorize the app via quire's auth page. It will automatically redirect you to your localhost to validate the auth tokens

## Commands
There different commands available for the quire cli

### `checkout`
This helps you to checkout a repository with the title of a task you choose and will eventually update your tasks status

For this command to run successfully, you have to provide it a project id which is usually found within the quire url:
```
https://quire.io/w/YOUR_PROJECT_ID
```
You can also change this ID by going to `https://quire.io/w/YOUR_PROJECT_ID?view=setting` and changing the *Project URL*

## Contribution
You're definitely welcome to contribute to this project. Just do the usual forking and making a PR for me :)

## New Ideas?
Just create a PR for your idea and I'll take a look at it as soon as I can. 
