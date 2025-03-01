package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

go mod operation failed. This may mean that there are legitimate dependency issues with the "go.mod" definition in the repository and the updates performed by the gomod check. This branch can be cloned locally to debug the issue.

Command that caused error:
./godelw check compiles

Output:
Running compiles...
nobadfuncs/types.go:102:29: in.TypeParams undefined (type *types.TypeParamList has no field or method TypeParams)
Finished compiles
Check(s) produced output: [compiles]

*/
