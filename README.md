# Rename strings in files

sorry not had time to document this better..



# The problem..
- Packaging a react app `build` folder into a docker image sucks
- Building different variants of this into discrete different directories can only scale so far
- CI/CD stage needs to execute many (inefficient ) `npm run build` runs
- storing different folders that have tiny differences suck


# Solution
1. Run the build once
2. ensure all `REACT_APP_` env vars have consistent values.. EG:
  `REACT_APP_FOO` will have the default value of `DEFAULT_VALUE_FOO`

When booting the container
3. loop over all env vars (eg `REACT_APP_FOO`)
4. substitute the values inside the files with the env vars of the container


# blerb
This tool simply evaluates all ENV vars with the prefix of `REACT_APP_` and does some file rewriting..

For example you have the env var `REACT_APP_FOO` = `BAR` set
The script will:
- find all files in the subdir
- rename all strings in the files if they are called `DEFAULT_VALUE_FOO`


# TLDR
- when adding new `REACT_APP_X` variables to you files, ensure to build the assets with a default value of `DEFAULT_VALUE_X`
- why not a simple bash+sed script?  a static binary is better to distribute and safer than shell scripting + distributing shell environments into a container
