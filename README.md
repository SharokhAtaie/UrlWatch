# UrlWatch
`UrlWatch is a tool for watch content of urls.`


**Note:** Replace your credentials in `main.go` line 23, and put your urls in urls.txt file
```go
var (
	accessToken   = "ghp_xxxxxxxxxxxxxxxx"                 // Your Github access token
	UserName      = "admin"                                // Your Github Username
	email         = "admin@gmail.com"                      // Your email (optional)
	repoUrl       = "https://github.com/Username/Repo.git" // Your Github Repo
	repoDir       = "./Repo"                               // Local Directory to save data
	TelegramToken = "123456:xxxxxxxxxxxxx"                 // Your Telegram bot token
	ChatID        = int64(-123456789)                      // Your ChatID of telegram bot
)
```


## Installation
```
git clone https://github.com/SharokhAtaie/UrlWatch.git
cd UrlWatch
go mod tidy
```

## Usage

##### If you want to run like cron job in linux, you can use this tool:
[hakcron](https://github.com/hakluke/hakcron)
```
go install github.com/hakluke/hakcron@latest

hakcron -f "daily" -c "go run main.go" // this run UrlWatch tool everyday
```

### Thanks :)
