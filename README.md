# Jmxbeat

Welcome to Jmxbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/radoondas`

## Getting Started with Jmxbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Init Project
To get running with Jmxbeat, run the following command:

```
make init
```

To commit the first version before you modify it, run:

```
make commit
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Jmxbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/radoondas/jmxbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Jmxbeat run the command below. This will generate a binary
in the same directory with the name jmxbeat.

```
make
```


### Run

To run Jmxbeat with debugging output enabled, run:

```
./jmxbeat -c jmxbeat.yml -e -d "*"
```


### Test

To test Jmxbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


### Package

To be able to package Jmxbeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Jmxbeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/jmxbeat.template.json and etc/jmxbeat.asciidoc

```
make update
```


### Cleanup

To clean  Jmxbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Jmxbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/radoondas
cd ${GOPATH}/github.com/radoondas
git clone https://github.com/radoondas/jmxbeat
```
