# tank

`tank` is a `game` that allows `anyone` to `have a bit of fun and take a break`.

You are in a `tank` at a fixed location and the `target` (that you want to shoot) is moving towards you at a fixed rate. You `choose the shot angle` to launch projectiles at the target, and the system will tell if you have `hit the target or by how much you missed the target`. This continues `until you either hit the target or it crushes you!`

You have the option of being right in the middle of the action with `manual shot mode`, or you can put tank in `auto shot mode` and watch the battle manager do its thing and narrow down on the target.

Either way, enjoy taking a break and trying to to get crushed by the oncoming target!!

## Prerequisites

Not much here. You can run tank on `Windows, Linux, or Mac`.

If you would like to compile, test, build the code then you will need `Golang` installed.

## Installing tank

To install tank, follow these steps:

Linux and macOS:
```
copy "tank" to your machine
chmod 775 tank # to make it executable
```

Windows:
```
copy "tank.exe" to your machine
```
## Using tank

To use tank, follow these steps:

Linux and macOS:
```
cd <dir that contains tank>
./tank
```

Windows:
```
cd <dir that contains tank.exe>
tank.exe
```

### Game Options:
```
Usage of ./tank:
  -a	Auto Shoot Mode (default - Single Player)
  -d float
    	Detonation Radius (meters) (default 20)
  -e	English Units (default - Metric)
  -p	Pause Target During Shot Decision (default - Realtime Target Movement)
```

#### Display

At startup, tank will tell you how it is configured:
```
Play Mode: Single Player
Play Mode: Realtime Target Movement
Units: Metric
Detonation Radius = 20.0 meters
```

You will be presented with a text display for each shot that tells you what the current situation is:
```
==================================
Projectile Velocity  = 382.6 meters/sec
Max Projectile Range = 14936.2 meters
Target Velocity      = 17.7 kilometers/hour
Target Velocity      = 4.9 meters/sec
Current Target Range = 12015.9 meters
----------------------------------
Enter a shot angle from 1.0 to 45.0 degrees:
```

After each shot, the resulting information will be displayed:
```
Taking shot #1 at 22.00 degrees.
Shot #1 took 29.2 seconds, and went 4192.0 meters.
Target Range = 11852.1 meters at time of impact.
<< Undershot target by -7660.1 meters.

 /~~~~~~~~~~~\
/--------+----\----+---------+---------T----------|
```
#### Timeline Explained

##### Flight Path

```
 /~~~~~~~~~~~\
```
The projectile will fly above the Target Path and stop at the end of the Flight Path giving a visual representation of how long the projectile is in the air.

##### Target Path

```
/--------+----\----+---------+---------T----------|
```

The Target Path is fixed and represents the maximum range of the projectile. There are several different symbols that appear on the Target Path at different times:
```
"/" - This is where you are, at a fixed point at the beginning of the Target Path.
"+" - This is a divider marker, breaking the Target Path into even chunks.
"|" - This is the end of the Target Path, the maximum range of the projectile.
"T" - This is the Target as it approaches you.
"\" - This represents a "miss" of the Target and you will see how close to the Target you are.
"*" - This represents a "hit" of the Target, it has been destroyed.
```

### Shot Options Explained

#### Manual Shot Mode
By default, tank runs in manual shot mode (you have to select the target angle) with real-time target movement (the target is moving all the time, while you are thinking what shot angle to use).

By selecting the `-p` option, the tank will only move when you `take a shot` - this gives you time to consider your shot angle with no pressure.

#### Auto Shot Mode
By selecting the `-a` option, the game will `play itself` and you can just sit back, relax, and see what happens. This can be extremely satisfying as you watch the drama unfold before you.

## Contributing to tank
To contribute to tank, follow these steps:

1. Fork this repository.
2. Create a branch: `git checkout -b <branch_name>`.
3. Make your changes and commit them: `git commit -m '<commit_message>'`
4. Push to the original branch: `git push origin tank/<location>`
5. Create the pull request.

Alternatively see the GitHub documentation on [creating a pull request](https://help.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request).


## License

This project uses the following license: [MIT License](https://github.com/scottballenger/tank/blob/main/LICENSE).