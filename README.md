Introduction
============
Packagebeat is a [Beat](https://www.elastic.co/products/beats)
for collecting information about system packages
from package managers and shipping it to [Elasticsearch](https://www.elastic.co/products/elasticsearch)

Package Managers
=====================

Packagebeat currently supports the following package managers:

 * dpkg (Debian, Ubuntu)
 * RPM (Fedora, CentOS, RHEL)

with hopes for supporting the following
(please consider contributing any of these):

 * [pip](https://pip.pypa.io/) (Python)
 * [gem](http://guides.rubygems.org/command-reference/#gem-list) (Ruby)
 * [npm](https://www.npmjs.com/) (node.js)
 * [chocolatey](https://chocolatey.org/) (Windows)
 * pacman (ArchLinux)
 * [nix](https://nixos.org/nix/) (NixOS)
 * [guix](https://www.gnu.org/software/guix/) (GuixSD)

Download
==========

Binaries are available on the [releases page](https://github.com/joehillen/packagebeat/releases).

Install
=========

The release package contains the following:

 * `packagebeat` binary
 * Example `packagebeat.yml`
 * The Elasticsearch mapping template: `packagebeat.template.json`

Install the mapping template before running Packagebeat:
```
curl -XPUT 'http://localhost:9200/_template/packagebeat' -d@packagebeat.template.json
```

Data
=====

Package information data is stored in the following format:

```json
{
   "@timestamp": "2099-01-01T00:00:00.000Z",
   "beat": {
     "hostname": "863bc3d673ad",
     "name": "863bc3d673ad"
   },
   "type": "package",
   "manager": "dpkg",
   "name": "tar",
   "version": "1.27.1-2+b1",
   "summary": "GNU version of the tar archiving utility",
   "architecture": "amd64"
 }
```

Building
==========

```
go get github.com/joehillen/packagebeat
```

Testing
=========

**Unit Tests:**

```
go test ./...
```

**Integration Testing:**

Testing on different Linux distributions is done using [docker-compose](https://docs.docker.com/compose/):

```
docker-compose up
```

You can inspect the results using Kibana at http://localhost:5601
