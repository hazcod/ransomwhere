# Ransomwhere

A Proof of Concept ransomware sample that encrypts your files to test out your ransomware detection & prevention strategies.</br>
I am not responsible for any damage caused by this software.

## Building

```shell
# with make and Go installed
% make build
```

## Usage

```shell
# straight from source
% make FLAGS="-log=warn -delete=false"

# from the binary
% ./app -log=warn -delete=false
```