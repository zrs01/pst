# Program Specification Tools
Tool to generate program specification document in .docx format

## Usage

```sh
NAME:
   dst - Database schema tool

USAGE:
   pst [global options] command [command options] [arguments...]

VERSION:
   development

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Config file
   --debug, -d               Debug mode (default: false)
   --help, -h                show help (default: false)
   --input value, -i value   Input file
   --output value, -o value  Output file
   --version, -v             print the version (default: false)
```


### Example

```sh
# -- Definition file
# example/example.yml


# -- Generate
$ pst -i sample.yml -o sample.docx
```

## Configuration

You may create configuration file to custom the properties of the output

Create a file `config.yml` with below content

```yml
# font name, e.g. Calibre. Default: Arial
fontfamily: Arial
# font size. Default: 10
fontsize: 10
logging:
  # available level: PANIC, FATAL, ERROR, WARN, INFO, DEBUG, TRACE. Default: INFO
  level: INFO
```

Pass the configuration file to the application when execute, e.g.

```sh
$ pst -i sample.yml -o sample.docx -c config.yml
```