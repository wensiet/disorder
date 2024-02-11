## Disorder

Project, that allow s you to store some files in discord,
however it is not recommended to do it, because discord will
probably purge your data.

## How to use

1. Get discord app and bot credentials
2. Create `db/files.db` in the root of the project
3. Fill the `config/config.yaml` according to the template
4. Build the program `go build main.go`
5. Run the program `./main`
6. Launch frontend, `build/index.html`

## How it works

As far as discord allows to store files in messages, we can
use channels as a storage. For a regular user there is a
limit of 25MB per file, but for a premium user it is 100MB. So by default
it is recommended to set bucket.size to 25MB. So, we need to split
the file into several chunks, discord CDN is public, so we will need
to cypher the files, I used AES-128 for that. The key will be stored at:
`/keys/aes`, it is highly recommended to have its copy in another place. Also, to
ensure that the data is not corrupted, we will store the checksum of each chunk and the
whole file itself. In the SQLite database we will store the metadata of the file, such as
the name, the size, the amount of chunks, chunks metadata and so on.