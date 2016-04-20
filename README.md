# Jmxproxybeat

**Welcome to Jmxproxybeat - simple GO JMX client.**

This beat request JMX metrics from running JVM of Apache Tomcat and sends them to Logstash or Elasticsearch.
Because Jmxproxybeat is not using JAVA, it is lightweight on the system. JMX metrics are requested through 'JMX_Proxy_servlet' configured in Tomcat for HTTP listener. The JMX Proxy Servlet is a lightweight proxy to get and set the tomcat internals.   

http://webserver/manager/jmxproxy/?get=BEANNAME&att=MYATTRIBUTE&key=MYKEY

## Tomcat configuration
Tomcat [JMX Proxy servlet](https://tomcat.apache.org/tomcat-8.0-doc/manager-howto.html#Using_the_JMX_Proxy_Servlet) configuration


## Getting Started with Jmxproxybeat

Ensure that this folder is at the following location:
`${GOPATH}/github.com/radoondas`

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Init Project
To get running with Jmxproxybeat, run the following command:

```
make init
```

To commit the first version before you modify it, run:

```
make commit
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Jmxproxybeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/radoondas/jmxproxybeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Jmxproxybeat run the command below. This will generate a binary
in the same directory with the name jmxproxybeat.

```
make
```


### Run

To run Jmxproxybeat with debugging output enabled, run:

```
./jmxproxybeat -c jmxproxybeat.yml -e -d "*"
```


### Test

To test Jmxproxybeat, run the following command:

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

To be able to package Jmxproxybeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Jmxproxybeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/jmxproxybeat.template.json and etc/jmxproxybeat.asciidoc

```
make update
```


### Cleanup

To clean  Jmxproxybeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Jmxproxybeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/radoondas
cd ${GOPATH}/github.com/radoondas
git clone https://github.com/radoondas/jmxproxybeat
```
