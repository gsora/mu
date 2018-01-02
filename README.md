# `mu`: a stupidly blind URL shortener

You'll never find a stupider URL shortener, I guarantee you!

## Why!

Why not? 

## How does it work?

Easy! This repository is `go get`-able, no external dependencies are needed.

Then, just run it!

Shortened links are saved on-disk in the `linkdb` file found in the folder where `mu` resides.

`mu` **does not** implements HTTPS support, just use a reverse proxy (nginx, Apache, or [Caddy](https://caddyserver.com/)).

## Endpoints

**Endpoint**|**Method**|**Arguments**
:-----:|:-----:|:-----:
|`/add`|**`GET`**|`url`, URL being shortened|

## Bonus, cloc!

```
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               2             42             23            149
-------------------------------------------------------------------------------
SUM:                             2             42             23            149
-------------------------------------------------------------------------------
```
