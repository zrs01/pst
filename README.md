# Program Specification Tools
Tool to generate program specification document from .yml to .docx

Note: the string in the .yml (except `module/name` and `scenarios/desc`) allow input string or string array.

e.g.
```yml
...
# single string
  name: this is single line of string
# array of string
  name: [first line, second line]
  # or
  name:
    - first line
    - second line
...
```


## Usage

```sh
NAME:
   pst - Program specfication tool

USAGE:
   pst [global options] command [command options] [arguments...]

VERSION:
   development

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value    config file
   --debug, -d                 debug mode (default: false)
   --document value, -m value  existing .docx file
   --help, -h                  show help (default: false)
   --input value, -i value     input file
   --output value, -o value    output file
   --version, -v               print the version (default: false)
```


### Example

```sh
# -- Sample file can be found at example/example.yml

# -- Generate in new document
$ pst -i sample.yml -o sample.docx

# -- Append to existing document
$ pst -i sample.yml -o sample.docx -m spec.docx

# -- Support multiple files
$ pst -i sample1.yml,sample2.yml,sample3.yml -o sample.docx

# -- Support wildcard input files
$ pst -i samp*.yml -o sample.docx
```


## Configuration

You may create configuration file to custom the properties of the output

Create a file `config.yml` with below content

```yml
# font name, e.g. Calibri. Default: Arial
fontfamily: Calibri
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