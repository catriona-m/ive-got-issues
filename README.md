# ivegotissues

A command-line tool for dealing with github issues.

## Commands

### addLabels
Adds labels to issues that match the filtering criteria

### addComment
Adds a comment to issues that match the filtering criteria

## Examples

Add a label to any issue where the body of text matches the regex passed to --content-matches flag:
```
ivegotissues addLabels --gh-owner "catriona-m" --gh-repo "ivegotissues" --labels "v2" --content-matches "Version=(|v)2\..[0-9]*\..[0-9]*" --state "open"
```

Add a comment to any issue that has the input labels attached to it:

```
ivegotissues addComment --gh-owner "catriona-m" --gh-repo "ivegotissues" --labels-filter "v2" --state "open" --comment "It looks like you are using a legacy version, please consider upgrading to 3.x.x"
```

## Notes

- The --dry-run flag is `true` by default - to make actual changes to issues you need to explicitly set --dry-run=false
- A Github access token is required to make the requests and is set via the environment variable `IGI_GITHUB_TOKEN`
