# Debian Utils

Package provide some debian utils

## Install

`go get github.com/bborbe/debian_utils/bin/create_debian_package`

`go get github.com/bborbe/debian_utils/bin/extract_zip`

`go get github.com/bborbe/debian_utils/bin/update_apt_source_list`

## Create Debian Package

```
create_debian_package \
-loglevel=DEBUG \
-version=1.2.3 \
-config=create_debian_package_config.json
```

### Sample Config

```
{
  "name": "debian-utils",
  "section": "utils",
  "priority": "optional",
  "architecture": "amd64",
  "maintainer": "Benjamin Borbe <bborbe@rocketnews.de>",
  "description": "Debian Package Utils",
  "postinst": "src/github.com/bborbe/debian_utils/postinst",
  "postrm": "src/github.com/bborbe/debian_utils/postrm",
  "preinst": "src/github.com/bborbe/debian_utils/preinst",
  "prerm": "src/github.com/bborbe/debian_utils/prerm",
  "files": [
    {
      "source": "bin/update_apt_source_list",
      "target": "/opt/debian_utils/bin/update_apt_source_list"
    },
    {
      "source": "bin/create_debian_package",
      "target": "/opt/debian_utils/bin/create_debian_package"
    }
  ]
}
```

## Update Apt-Repo

```
update_apt_source_list \
-loglevel=DEBUG \
-path /etc/apt/sources.list.d/aptly-unstable.benjamin-borbe.de.list
```

## Documentation

[GoDoc](http://godoc.org/github.com/bborbe/debian_utils/)

## Continuous integration

[Jenkins](https://www.benjamin-borbe.de/jenkins/job/Go-Debian-Utils/)

## Copyright and license

    Copyright (c) 2016, Benjamin Borbe <bborbe@rocketnews.de>
    All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions are
    met:
    
       * Redistributions of source code must retain the above copyright
         notice, this list of conditions and the following disclaimer.
       * Redistributions in binary form must reproduce the above
         copyright notice, this list of conditions and the following
         disclaimer in the documentation and/or other materials provided
         with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
