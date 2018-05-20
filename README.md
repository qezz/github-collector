# Collect data about users from github

## How to

1. Receive your personal token: https://github.com/settings/tokens
   * Creeate new if needed
   * set access to **user**
   * Generate token


## Why use API token?

* GitHub imposes a rate limit on all API clients. Unauthenticated
  clients are limited to 60 requests per hour, while authenticated
  clients can make up to 5,000 requests per hour. 
* actually 32 per minute == 1920 per hour

## Deps

```
# opencv
go get -u -d gocv.io/x/gocv
```

## Run on macOS

```
brew cleanup opencv # ensures that you have only one version installed

source ./env.sh
```



