Telegram Secret Santa bot
-------------------------

**Build**

```bash
go build -o ./secretsantabot ./cmd/secretsantabot/main.go
```

**Start bot**
```bash
TELEGRAM_TOKEN=000000:AAAAAAAAAA DATABASE_URL=postgresql://user:password@127.0.0.1:5432/db ./secretsantabot
```

***

**Commands**
```
/enroll - enroll the game
/disenroll - stop your enroll (only before magic starts)
/list - list all enrolling people
/magic - start the game (only admin)
/my - resend magic info for you (only in private chat wi me)
/help - show this message
```

***

**Heroku create new app**

configure manifest module

```bash
$ heroku update beta
$ heroku plugins:install @heroku-cli/plugin-manifest
```

create new app  

```bash
$ heroku apps:create --manifest -n -s container --region eu APP_NAME
```