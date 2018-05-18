# Debian Utils

Package provide some debian utils

## Create Debian Package

create_debian_package create a Debian-Package with instructions out of a Json config.  

`go get github.com/seibert-media/debian-utils/cmd/create_debian_package`

```
create_debian_package \
-logtostderr \
-v=2 \
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
  "postinst": "src/github.com/seibert-media/debian-utils/postinst",
  "postrm": "src/github.com/seibert-media/debian-utils/postrm",
  "preinst": "src/github.com/seibert-media/debian-utils/preinst",
  "prerm": "src/github.com/seibert-media/debian-utils/prerm",
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

update_available_apt_source_list checks if a debian repo has changed and update_apt_source_list fetches the new sources.

`go get github.com/seibert-media/debian-utils/cmd/update_available_apt_source_list`

`go get github.com/seibert-media/debian-utils/cmd/update_apt_source_list`

```
update_available_apt_source_list \
-logtostderr \
-v=2 \
-path /etc/apt/sources.list.d/aptly-unstable.benjamin-borbe.de.list
```

```
update_apt_source_list \
-logtostderr \
-v=2 \
-path /etc/apt/sources.list.d/aptly-unstable.benjamin-borbe.de.list
```

## Extract Zip

`go get github.com/seibert-media/debian-utils/cmd/extract_zip`
