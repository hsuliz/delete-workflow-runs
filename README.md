# delete-workflow-runs

Deletes all of workflow runs in given repository.

## Description

GitHub actions do not have single button to delete all runs for workflow.

## How to use

You should generate GitHub personal access token (classic) \
https://github.com/settings/tokens

### Build

```shell
make build
```

### Run

By default deletes all workflow runs in given repository:

```shell
./dltwfrns -o <GITHUB_USER> -r <GITHUB_REPO> -bt <GITHUB_TOKEN>
```

To delete from specific workflow use `-wn` flag:

```shell
./dltwfrns -o <GITHUB_USER> -r <GITHUB_REPO> -bt <GITHUB_TOKEN> -wn <WORKFLOW_NAME>
```