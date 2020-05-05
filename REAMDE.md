## Sparse CRC32

```bash
go get github.com/furstenheim/sparse_crc32
```

This modules exports a function to compute crc32 on sparse files. If you want to compute a CRC32 on a file of 100GB which is mostly empty, then this is your library.

```go
// crc32 on a 4GB file with two non null bytes
file := sparse_crc32.SparseFile{[]sparse_crc32.SparseByte{{0, 1}, {4, 55}}, 1 << 32},
crc32 := sparse_crc32.IEEESparse(file)
```
