# autopush

A solution to when we forget to commit and push changes on a remote repo. This tool will
automatically push all changes to the remote of the working directory where the tool was
executed. The commit will have the commit message

```txt
automated commit by autopush
```

> autopush probably should *not be used* when working in a large team, as
> constant commits can lead to unexpected issues and conflicts. autopush
> is intended for smaller teams, ideally a single developer who forgets to
> commit their changes to their remote repo.

## Installation

### Unix

To use autopush, first make sure you have Go installed.

```sh
go -v
```

> If you do not have Go installed, checkout the official docs for installation
> instructions.

1. Clone the repo

```sh
git clone https://github.com/villaleo/autopush.git
```

2. Change directory to the project directory

```sh
cd autopush/
```

3. Build the binary

```sh
go build -o bin/ main.go
```

4. Now that the binary is built, you can change directory to your project repo
and execute the binary in the background! In a new terminal window, run

```sh
cd my-cool-project
./../autopush/bin/main
```

> The `./` is used to execute a binary. Everything after the `./` is the path
> to the binary you wish to execute.

### macOS

> If you are running macOS, then you may **not** need to build the Go program.
> You may be able to just execute the binary in `bin/main`.

```sh
cd my-cool-project
./../autopush/bin/main
```

## Contributing

Feel free to create any issues for feature requests and/or bug reports. Thanks!
