# autopush

A solution to when we forget to commit and push changes on a remote repo. This tool will
automatically push all changes to the remote of the working directory where the tool was
executed. 

> autopush probably should *not be used* when working in a large team, as
> constant commits can lead to unexpected issues and conflicts. autopush
> is intended for new developers commiting to a linear repo.

## Installation

1. Clone the repo

```sh
git clone https://github.com/villaleo/autopush.git
```

2. Change directory to the project directory

```sh
cd autopush/
```

3. Run the Makefile to install all dependencies. This will ask for your password
to add the compiled binary to your PATH.

```sh
make build
```

## Contributing

Feel free to create any issues for feature requests and/or bug reports. Thanks!
