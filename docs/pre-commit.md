# Pre commit hooks

This repository uses pre-commit hooks to ensure code quality and consistency. The hooks are defined in the `.pre-commit-config.yaml` file.


## Please install pre-commits before working in this repository
```shell
pip install --upgrade pre-commit
pre-commit install --install-hooks
```

## Auto update all pre-hooks
```shell
pre-commit autoupdate
```

## Run pre-commit hooks on all files
```shell
pre-commit run --all-files
```

By default, the pre-commit hooks will run on all files. If you want to run the hooks on specific files, you can do so by specifying the file names as arguments to the `pre-commit run` command.
