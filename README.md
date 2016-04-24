# Jmxproxybeat

**Welcome to Jmxproxybeat - simple beat for JMXProxyServlet for Apache Tomcat to retrieve JMX metrics.**

This is still **development version** and I expect changes based on feedback.

This beat retrieves JMX metrics from running JVM of Apache Tomcat and sends them to Logstash or Elasticsearch.
JMX metrics are requested via 'JMX Proxy Servlet' configured and enabled in Tomcat for HTTP listener. JMX Proxy Servlet is a lightweight proxy to get and set the Tomcat internals.

Because Jmxproxybeat is not using JAVA, it is lightweight on the system and there is no need for JAVA to get JMX metrics.

## Tomcat configuration
General reading about Tomcat [JMX Proxy servlet](https://tomcat.apache.org/tomcat-8.0-doc/manager-howto.html#Using_the_JMX_Proxy_Servlet). 

In order to enable JMX Proxy Servlet in default Tomcat package, minimal configuration in **conf/tomcat-users.xml** is required. Tomcat restart is also required.
```xml
<role rolename="manager-jmx"/>
<user username="tomcat" password="s3cret" roles="manager-jmx"/>
```

Test if your Tomcat configuration works with your credentials.

Template scheme of the request of the Bean:
```
http://127.0.0.1:8080/manager/jmxproxy/?get=BEANNAME&att=MYATTRIBUTE&key=MYKEY
```

Example of request for **HeapMemoryUsage**
```
http://127.0.0.1:8080/manager/jmxproxy/?get=java.lang:type=Memory&att=HeapMemoryUsage&key=init
```

## Getting Started with Jmxproxybeat

Ensure that this folder is at the following location:
`${GOPATH}/github.com/radoondas`

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Build

To build the binary for Jmxproxybeat run the command below. This will generate a binary in the same directory with the name jmxproxybeat.

```
make
```


### Run

To run Jmxproxybeat with debugging output enabled, run:

```
./jmxproxybeat -c jmxproxybeat.yml -e -d "*"
```

### Example JSON output
```
{
  "_index": "jmxproxybeat-2016.04.20",
  "_type": "jmx",
  "_id": "AVQ0FOGeegQ15caFDGZ7",
  "_score": null,
  "_source": {
    "@timestamp": "2016-04-20T14:31:03.385Z",
    "bean": {
      "attribute": "MemoryUsed",
      "hostname": "127.0.0.1:8080",
      "name": "java.nio:type=BufferPool,name=direct",
      "value": 81920
    },
    "beat": {
      "hostname": "localhost",
      "name": "localhost"
    },
    "type": "jmx"
  }
```

### Test - not complete yet

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


### Package - not complete yet

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
