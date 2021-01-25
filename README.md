# tank

`tank` is a `game` that allows `anyone` to `have a bit of fun and take a break`.

You are in a `tank` at a fixed location and the `target` (that you want to shoot) is moving towards you at a fixed rate. You `choose the shot angle` to launch projectiles at the target, and the system will tell if you have `hit the target or by how much you missed the target`. This continues `until you either hit the target or it crushes you!`

You have the option of being right in the middle of the action with `manual shot mode`, or you can put tank in `auto shot mode` and watch the battle manager do its thing and narrow down on the target.

Either way, enjoy taking a break and trying to to get crushed by the oncoming target!!

## Prerequisites

Not much here. You can run tank on `Windows or Mac`.

If you would like to compile, test, and build the code then you will need `Golang` installed.

## Installing tank

To install tank, follow these steps:

macOS:
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

macOS:
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
  -a	Auto Shoot Mode (default - Manual Shot)
  -d float
    	Detonation Radius (meters) (default 20)
  -e	English Units (default - Metric)
  -m	Real-time Target Movement (default - Pause Target During Shot Decision)
  -p	Print Shot Profile
```
#### Random Values:

At startup, a random value is selected for `Projectile Velocity` (in the following range)
```
min Projectile Velocity = 300 meters/sec
max Projectile Velocity = 600 meters/sec
```
 and `Target Velocity` (in the following range)
```
min Target Velocity     = 0   kilometers/hour
max Target Velocity     = 60  kilometers/hour
```

#### Display

At startup, tank will tell you how it is configured:
```
Shoot Mode: Manual
Target Mode: Pause Target During Shot Decision
Units: Metric
Detonation Radius = 20.0 meters
```

You will be presented with a text display for each shot that tells you what the current situation is:
```
==================================
Projectile Velocity  = 558.5 meters/sec
Max Projectile Range = 31808.4 meters
Target Velocity      = 53.2 kilometers/hour
Target Velocity      = 14.8 meters/sec
Current Target Range = 26.6 kilometers
Current Target Range = 26640.1 meters
----------------------------------
Enter a shot angle from 1.0 to 45.0 degrees (0 to quit):
```

After each shot, the resulting information will be displayed:
```
Taking shot #1 at 22.00 degrees. Flight time is 42.7 seconds.
Shot #1 took 42.7 seconds, and went 22095.9 meters (22.1 kilometers).
Target Range = 26009.2 meters at time of impact.
<< Undershot target by -3913.3 meters.

 /~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\
/--------+---------+---------+----\----+T--------|
        6.4K     12.7K     19.1K     25.4K     31.8Kilometers
```
#### Timeline Explained

##### Flight Path

```
 /~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~\
```
The projectile will fly above the Target Path and stop at the end of the Flight Path giving a visual representation of how long the projectile is in the air.

##### Target Path

```
/--------+---------+---------+----\----+T--------|
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

##### Target Path Ruler/Legend

```
/--------+---------+---------+---------+---------|
        6.4K     12.7K     19.1K     25.4K     31.8Kilometers
```

The Ruler/Legend describes the distances that the Target Path represents. Each `+` has a number under it, in either `K`ilometers or `M`iles. The last distance spells out the full units, so as to not clutter the entire timeline.

### Shot Options Explained

#### Manual Shot Mode
By default, tank runs in manual shot mode (you have to select the target angle) with no target movement. The tank will only move when you `take a shot` - this gives you time to consider your shot angle with no pressure. 

By selecting the `-m` option, the target will move all the time - while you are thinking what shot angle to use. This represents a challenge to make your decisions while impending doom is looming!

#### Auto Shot Mode
By selecting the `-a` option, the game will `play itself` and you can just sit back, relax, and see what happens. This can be extremely satisfying as you watch the drama unfold before you.

#### Print Shot Profile

Selecting the `-p` option will print out a table of the shot angles from 1-45 with their corresponding ranges and times for the (random)`Projectile Velocity` in your run. 

You can use this information to decide what shot angles to take. Here's what the Shot Profile looks like for a Projectile Velocity = 567.4 meters/sec:

```
Shot Profile:
+-------+------------+-------+
| Angle | Shot Range | Time  |
| (deg) |   (meters) | (sec) |
+-------+------------+-------+
|   1.0 |     1145.7 |   2.0 |
|   2.0 |     2290.0 |   4.0 |
|   3.0 |     3431.4 |   6.1 |
|   4.0 |     4568.8 |   8.1 |
|   5.0 |     5700.5 |  10.1 |
     :            :       :
|  40.0 |    32329.1 |  74.4 |
|  41.0 |    32508.4 |  75.9 |
|  42.0 |    32648.0 |  77.4 |
|  43.0 |    32747.9 |  78.9 |
|  44.0 |    32807.9 |  80.4 |
|  45.0 |    32827.9 |  81.8 |
+-------+------------+-------+
```

## Building/Testing tank
`tank` is developed in Golang. You will need to download Golang from https://golang.org/doc/install. You can install additional developer tools such as an IDE if you would like, but it is not required.

### Golang Version
This code was compiled with `go version go1.15.7 darwin/amd64`. Run `go version` to see what you are using.

### Compile the Code and Build Executables

To build the code and create the stand-alone executable for your platform, just run the following command:

```
cd tank
go build
```

macOS:
This will create the executable `tank` that you can run.

Windows:
This will create the executable `tank.exe` that you can run.

#### Compiling the Code for other Platforms

For the complete list of operating systems and architectures that can be cross compiled, see https://golang.org/doc/install/source#environment

##### Compiling for Windows from macOS

If you are on a macOS platform and want to create an executable for Windows, then you would run the following:

```
cd tank
GOOS=windows go build
```

This will create the executable `tank.exe` that you can run on Windows.

##### Compiling for macOS from Windows

If you are on a Windows platform and want to create an executable for macOS, then you would run the following:

```
cd tank
GOOS=darwin go build
```

This will create the executable `tank` that you can run on macOS.

### Run Unit Tests

To run the unit tests for your platform, just run the following command:

```
cd tank
go test
```

Upon execution, you should see something that ends with:
```
PASS
ok      tank    0.319s
```

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