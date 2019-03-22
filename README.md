# awesome-wallpaper

Automatically change wallpaper by crontab-like schedule syntax.

| Arguments | Default | Description |
|----------|-------------|---|
| --schedule | * * * * * | Optional. Crontab-like syntax schedule |
| --keywords | (empty) | Optional. Keywords to search, separated by commas|

## Example

Set random wallpaper every minute:

```
awesome-wallpaper
```

Set random wallpaper at 7:00 AM every day:

```
awesome-wallpaper --schedule=0 7 * * *
```

With keywords:

```
awesome-wallpaper --schedule=0 7 * * * --keyword=cat,monster
```

or

```
awesome-wallpaper --keyword=cat,monster
```