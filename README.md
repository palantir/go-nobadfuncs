<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/go-nobadfuncs"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

go-nobadfuncs
=============
go-nobadfuncs verifies that a set of specified functions are not referenced in the packages being checked. It can be used to deny specific functions that should not typically be referenced or called. It is possible to explicitly allow uses of deny-listed functions by adding a comment to the line before the calling line.

Usage
-----
go-nobadfuncs takes the path to the packages that should be checked for function calls. It also takes configuration (as JSON) that specifies the blacklisted functions.

The function signatures that are blacklisted are full function signatures consisting of the fully qualified package name or receiver, name, parameter types and return types. Examples:

```
func (*net/http.Client).Do(*net/http.Request) (*net/http.Response, error)
func fmt.Println(...interface{}) (int, error)
```

go-nobadfuncs can be run with the following flags:

* `--print-all` flag to print all of the function references in the provided packages. The output can be used as the basis for determining the signatures for blacklist functions.
* `--config-json` flag to run with the JSON configuration for the check
