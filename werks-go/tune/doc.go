/*
Package tune provides a means of incorporating tunable constants into an
application where the defaults are loaded with the executable.  Specific,
"tuned" values can then be imported on top of the defaults from an alternative
source (json/yaml).

Look at the mapstructure package for inspiration.
  - https://pkg.go.dev/github.com/mitchellh/mapstructure ^^^ isn't what I was
    looking for, there's another package that eloquently provides access to
    values where the type of the value isn't known until run-time.  .. keep
    looking
  - pflag package <-- is where the inspiration came from.
*/
package tune
