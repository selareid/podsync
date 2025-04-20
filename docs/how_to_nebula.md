# How to get token
1. Log in to nebula.tv
2. Open Network tab in developer console
3. Refresh page
5. Find request to `/api/v1/authorization`
6. Under request headers find `Authorization Token: <token>`
7. Copy the value of `<token>` and use in `config.toml`

# How to get nebula feed
Setup your config.toml as below.

For each feed, you must use url format `https://nebula.tv/<channel name>`. Only channels are supported at this time.

An example config.toml for Nebula:
```
[server]
port = 8080
hostname = "http://localhost:8080"

[storage]
  [storage.local]
  data_dir = "./app/data/"

[tokens]
nebula = "<your token>"

[feeds]
    [feeds.jetlag]
    youtube_dl_args = ["--embed-subs"]
    url = "https://nebula.tv/jetlag"
    max_height = 1080
    page_size = 3
    opml = true
```

I like to use feed option
`youtube_dl_args = ["--embed-subs"]`
to get subtitles for jetlag but that's up to you.