# How to get nebula feed
Setup your config.toml as below.

Under tokens put `nebula = "nil"`.
For each feed, include
`youtube_dl_args = ["--username", "<username>", "--password", "<password>", "--embed-subs"]`.
You must use url format `https://nebula.tv/<channel name>`. Only channels are supported at this time.

An example config.toml for Nebula:
```
[server]
port = 8080
hostname = "http://localhost:8080"

[storage]
  [storage.local]
  data_dir = "./app/data/"

[tokens]
nebula = "nil"

[feeds]
    [feeds.wendover]
    youtube_dl_args = ["--username", "<redacted>", "--password", "<redacted>", "--embed-subs"]
    url = "https://nebula.tv/wendover"
    max_height = 1080
    page_size = 4
    opml = true
```