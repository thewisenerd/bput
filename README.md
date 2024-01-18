bput
====

_bput_ is a command line tool to upload blobs to Azure Blob Storage.

_bput_ assumes that you have
a [System-assigned Managed Identity](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/overview)
assigned on your Azure VM,
and the identity has been provided Storage Blob Data Contributor Role to the Azure Blob Storage
account that you want to upload to.

_bput_ uses the official [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go).

# Usage

```
Usage: bput [options] <file1> <file2>...
  -a string
    	account
  -b string
    	bucket/container
  -p string
    	path prefix

Example: bput -a myaccount -b mycontainer file1 file1.md5sum
```

# Compilation

```
go mod download
make
```
