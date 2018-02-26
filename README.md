# Backup

Backup is a tool for backing up data to a remote server.

It consists of two parts, a server, and a client.

On the client side, backup uses gpg to encrypt local archives, and the server doesn't need to know or
care about the contents of the file, or the contents of the archive.

Archiving is done via tar, and this means that even without backup, you can still access your files
if you only have access to your gpg key and the files themselves.

Backup aims to deduplicate archives, and thus provide incremental backups. It transfers the archives
between client and server via TLS secured HTTP.

Backup is a work in progress, and may never achieve all of these goals, as it is mostly a playground
experiment for the author to understand more about Go.

---

Right now, backup creates a tar.gz and an md5 of that archive. Nothing more.
