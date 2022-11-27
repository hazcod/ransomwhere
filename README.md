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
# straight from source
% make FLAGS="-log=warn -delete=false -mode=encrypt"

# from the binary
% ./app -log=warn -delete=false -mode=encrypt

# encrypt, delete original files and wipe backups like a real ransomware (DANGEROUS)
% ./app -delete=true -wipe=true

# revert the ransom operation and restore any files
% ./app -mode=decrypt
```
