# ive-got-issues

This is a command-line tool for dealing with GitHub issues. This is intended to make sorting through issues in large repos faster and easier. 

It can automatically `label` or `comment` all issues in a specified repo that match the criteria input via the flags. These have a dry-run mode and the option to open a browser tab for each matching issue for users who prefer to double-check matching issues before commenting/labelling. 

It can `list` all issues in a specified repo that match the criteria input via the flags. This also has the option to list issues with `--linked-prs` (pull requests that have referenced the issues), which can help identify stale issues that are ready to close. `list` can also optionally open issues in browser tabs.

ive-got-issues is still a work in progress with more flags and commands to come in future!

## Installation

To install ive-got-issues from the command line, you can run:

`go install catriona-m/ive-got-issues`

## Commands

## label
Adds labels to issues that match the filtering criteria

### Examples

- Add a label to any issue where the body of the issue text matches the regex passed to --content flag:
```
ive-got-issues labels --gh-repo "catriona-m/ive-got-issues" --labels "v2" --content "Version=(|v)2\..[0-9]*\..[0-9]*" --state "open"
```

## comment
Adds a comment to issues that match the filtering criteria

### Examples

- Add a comment to any issue that has the input labels attached to it:
```
ive-got-issues comment --gh-repo "catriona-m/ive-got-issues" --labels "v2" --state "open" --comment "It looks like you are using a legacy version, please consider upgrading to v3"
```

### list
List all issues that match the filtering criteria

### Examples

- List all open issues with a `bug` label that are referenced in a merged pull request
```
ive-got-issues list --gh-repo "catriona-m/ive-got-issues" --labels "bug" --state "open" --pr-state "merged" --linked-prs=true
```

- List all open issues and the merged prs that are referenced in a merged pull request, open a browser tabs for each matching issue and linked pr in batches of 5 at a time
```
ive-got-issues list --gh-repo "catriona-m/ive-got-issues" --labels "bug" --state "open" --pr-state "merged" --linked-prs=true --open=true --open-prs=true --batch 5
```

- List all open issues with `bug` and `service/my-service` labels 
```
ive-got-issues list --gh-repo "catriona-m/ive-got-issues" --labels "bug,service/my-service" --state "open"
```

## Notes

- The --dry-run flag is `true` by default - to make actual changes to issues you need to explicitly set --dry-run=false
- A GitHub access token is required to make the requests and is set via the environment variable `IGI_GITHUB_TOKEN`
