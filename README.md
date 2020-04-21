# qGmail
A fast and lightweight CLI app which brings the power of [Gmail API](https://developers.google.com/gmail/api) to your terminal.

## Overview
qGmail, short for *query Gmail*, is a lightweight app written in [Go](https://golang.org). It aims to bring the powerful [Gmail API](https://developers.google.com/gmail/api) to your terminal in an easy-to-use manner so that you don't have to bother about the authentication stuff yourself.  
qGmail uses [Google's recommended authorization flow](https://developers.google.com/identity/protocols/oauth2/native-app), OAuth2(with PKCE extension), making its transactions with the API highly secure.

> In the initial releases, `qgmail` only fetches the number of unread mails but it is planned to include lots of other functionalities depending upon the feature requests. If you have one, feel free to write up a feature request through Issues.

### Supported Architectures
Currently, qGmail officially supports Linux distros only but it is planned to implement cross-compatibility across Windows, Mac and Linux.

## Choose How to Install
### 1. Binary Install
You can download pre-built binaries from [qGmail's releases](https://github.com/UtkarshVerma/qgmail/releases) page.
### 2. Building from Source
#### Prerequisite tools
* [Git](https://git-scm.com/)
* [Go](https://golang.org/dl)
#### Fetch from GitHub
qGmail uses the Go Modules support built into Go 1.11 to build. The easiest is to clone qGmail in a directory outside of `GOPATH`, as in the following example:
```bash
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/utkarshverma/qgmail.git
cd qgmail
go install
```
## Usage
To get started:
*  [Create OAuth2 Credentials](https://developers.google.com/identity/protocols/oauth2#1-obtain-oauth-20-credentials-from-the-google-api-console) using Google API Console.
* Make sure your API Console client has access to Gmail API.
* Download your OAuth2 credentials as `credentials.json` and paste it in `~/.config/qgmail`.
* Run `qgmail init`.
* Follow the on-screen instructions to authorize qGmail.

Once the initialisation is done, just simply use `qgmail` to access the API.

### CLI
qGmail provides some CLI flags to help you configure it on the go. To know about them, just use `qgmail --help`, or `qgmail init --help`.

## Why should you trust qGmail?
**qGmail doesn't store your data in any  manner**, nor does it have access to your account. The client attains access to your account through an authorization token which Google grants to **your PC** after you accept the user consent. Since the authorization token stays only on your PC and qGmail doesn't have any functionality to upload files, so rest assured as **your account's access stays with you only**!

## Contributing to qGmail
qGmail welcomes all sorts of contributions. If you want to contribute to qGmail, there are three ways you can do so:
### 1. Feature Requests
If you feel that there's a feature this app is lacking, [post a feature request](https://github.com/UtkarshVerma/qgmail/issues/new).
### 2. Testing
There are cases when fixing a bug triggers another bug, which I might not notice during testing myself. Having more users would be a great help to close in on such bugs.
If you encounter some bugs, [post an issue regarding the bug](https://github.com/UtkarshVerma/qgmail/issues/new).
### 3. Pull Requests
I'm maintaining this open-source project as a personal interest without any profits, hence I'll only be developing qGmail if I have the time and motivation to do so. Therefore, PRs from other developers are highly welcomed.

## Dependencies
qGmail stands on the shoulder of many great open-source libraries, in lexical order:
| **Dependency**|**License**|
|:---:|:---:|
|[google.golang.org/api](https://google.golang.org/api)|BSD-3 Clause|
|[github.com/mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)|MIT License|
|[golang.org/x/net](https://golang.org/x/net)| BSD-3 Clause|
|[golang.org/x/oauth2](https://golang.org/x/oauth2)| BSD-3 Clause|
