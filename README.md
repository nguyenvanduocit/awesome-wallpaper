# awesome-wallpaper

Automatically change wallpaper by crontab-like schedule syntax.

## Install

```
go get github.com/nguyenvanduocit/awesome-wallpaper
```

## Arguments

| Arguments  | Default   | Type   | Description |
|------------|-----------|--------|-------------|
| --schedule | * * * * * | String | Optional. Crontab-like syntax schedule |
| --keywords | (empty)   | String | Optional. Keywords to search, separated by commas|

## Example 

Set random wallpaper every minute:

```
awesome-wallpaper
```

Set random wallpaper at 7:00 AM every day:

```
awesome-wallpaper --schedule="0 7 * * *"
```

With keywords:

```
awesome-wallpaper --schedule="0 7 * * *" --keyword=cat,monster
```

or

```
awesome-wallpaper --keyword=cat,monster
```

## Crontab syntax

Here are the few quick references about crontab simple but powerful syntax.

```
*     *     *     *     *        

^     ^     ^     ^     ^
|     |     |     |     |
|     |     |     |     +----- day of week (0-6) (Sunday=0)
|     |     |     +------- month (1-12)
|     |     +--------- day of month (1-31)
|     +----------- hour (0-23)
+------------- min (0-59)
```

### Examples

+ `* * * * *` run on every minute
+ `10 * * * *` run at 0:10, 1:10 etc
+ `10 15 * * *` run at 15:10 every day
+ `* * 1 * *` run on every minute on 1st day of month
+ `0 0 1 1 *` Happy new year schedule
+ `0 0 * * 1` Run at midnight on every Monday

### Lists

+ `* 10,15,19 * * *` run at 10:00, 15:00 and 19:00
+ `1-15 * * * *` run at 1, 2, 3...15 minute of each hour
+ `0 0-5,10 * * *` run on every hour from 0-5 and in 10 oclock

### Steps
+ `*/2 * * * *` run every two minutes
+ `10 */3 * * *` run every 3 hours on 10th min
+ `0 12 */2 * *` run at noon on every two days
+ `1-59/2 * * * *` run every two minutes, but on odd minutes