REST Path and API
=================

`cfg` API
---------

This API serves up configuration files based on templates with substitutions
from the CMDB.  The configuration files may be in any format, and specificaly
are intended to be in the format suitable for consumption by the application
they are intended for.  This is a bridge from the CMDB's data to arbitrary
application configurations.

### `/cfg` API

```
GET /cfg                 - return a list of configuration file paths
GET /cfg/<prefix-path>   - list of configuration file paths under the prefix

POST /cfg + JSON         - declare a template path + admin operations

GET /cfg/<cfg-file-path> - return a populated configuration file (or 404)
PUT /cfg/<cfg-file-path> - update the template (or 404)
DEL /cfg/<cfg-file-path> - delete the template (or 404)
```

`cmdb` API
----------

The `cmdb` API provides operations against the CMDB.

_Note:  The REST paths below are lifted from the home-sensors project and should
be reassessed.  -EG 1/2024_

```
GET   /cmdb    ??

GET   /cmdb/ci   ?? list ci oids
POST  /cmdb/ci   create new ci, initialize config[0] with json in body

GET   /cmdb/ci/[oid]  return ci object, includes list of config versions,
                        but not config's themselves.
POST  /cmdb/ci/[oid]  commit the next config[new#] with json in body.

GET   /cmdb/ci/[oid]/[ver|tag] return specified config

POST  /cmdb/ci/search  -- search and return ci's matching search
                          criteria.

--

POST  /cmdb/mo/discover  -- discover an MO based on body-json criteria

GET   /cmdb/mo    return a list of MO oids

# POST  /cmdb/mo?template=oid  -- create a MO instance based on template[oid]
#                                 return new oid for MO

GET   /cmdb/mo/[oid]  return the MO for CI[oid]
PUT   /cmdb/mo/[oid]  initialize / reset the MO's state (json body)
POST  /cmdb/mo/[oid]  update state of the MO using json body
DEL   /cmdb/mo/[oid]  remove all state (aka 'unregister') the MO

# GET   /cmdb/mo/[oid]/status  get status
# POST  /cmdb/mo/[oid]/status  set status

```