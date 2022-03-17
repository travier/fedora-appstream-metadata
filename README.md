# Fedora Linux AppStream Metadata

[Operating System AppStream Metadata] for Fedora Linux.

Fedora Linux releases are fetched from
<https://admin.fedoraproject.org/pkgdb/api/collections/>.

Release and EOL dates are fetched from the schedule milestones at
<https://fedorapeople.org/groups/schedule/>.

Note that they do not exactly correspond to the date at
<https://docs.fedoraproject.org/en-US/releases/eol/>.

## Updating the metadata

The metadata file in this repo is auto-generated and updated using:

```
$ go run update-appstream-metadata.go
```

[Operating System AppStream Metadata]: https://www.freedesktop.org/software/appstream/docs/sect-Metadata-OS.html
