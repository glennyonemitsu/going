package main

var ExampleConfigGoing string = `
# This is the example config file with all configuration parameters for the 
# main going server.

# The location of going's pid file. Defaults to the directory of the loaded 
# config file.
#PidFile:

# The location of the program conf files. Defaults to the directory of the
# loaded config file.
#ProgramConfigDir:

# The location of the socket file. Defaults to the directory of the pid file.
#SocketFile:
`

const ExampleConfigProgram string = `
# This is the example config file with all configuration parameters for the 
# programs managed by the going server.

# Full path command to execute. 
# Required
#Command:

# Key: Value pairs for the command. If any Key: Value pairs are specified, you
# must also uncomment the parent "Environment" line as well.
# Optional
#Environment:
  #ExampleEnvVar: Foobar
  #ExampleSecretVar: This Is A Secret

# Easy to remember identifier used for internal and management commands. Each
# program must have a unique name.
# Required
#Name:

# OS level signal to use when shutting down this program through the going
# command interface.
# Optional
#StopSignal:

# Working directory for the command. Defaults to the working directory of the
# going server.
# Optional
#Dir:
`
