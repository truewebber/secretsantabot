Telegram Secret Santa bot
-------------------------

**Start bot**
```bash
ENV=prod LOGXI=* ./secret_hell_santa
```

You can set CONFIG_PATH to define config location. Default in the same path with bin
Config name it's ENV.yaml. Default ENV=dev

***

**Commands**
```
/enroll - enroll the game
/end - stop your enroll (only before magic starts)
/list - list all enrolling people
/magic - start the game (only admin)
/my - SecretHelSanta will resend magic info for you (only in private chat wi me)
/help - show this message
```

***

**Config example**
```yaml
token: 0000000:AAAAAAAAAAAAAAAAAAAAAAAAA
lock-on-chat-id: 0000000
admin-user-id: 0000000

rules:
  deny:
    1212: 2
    2323: 4
    -933: 6

```
 - for now only deny rules
 - token, lock-on-chat-id - required
 - admin-user-id - lock magic method
 - rules use TelegramID, you can see `user-id` in debug messages
 
