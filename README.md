# plfit-tablegen

A small utility to summarize `plfit` data.

## Install

```
go get ivo.qa/plfit-tablegen
cd plfit-tablegen
go build
```

## Prerequisites
 - [`plfit`](https://pypi.org/project/plfit/) must be installed
 - The .dat files must follow this naming convention regular expression:
 
 ```
 ^[\w\d\-]+_(?P<Frame>\d+)_(?P<Type>[\w\d\-]+)_(?P<Period>\d+)_[\w\d\-]+_(?P<StructNo>\d+)\.dat$
 ```


## Run

```
./plfit-tablegen <path containing .dat files> [--silent]
```

## Output

For example,

```
 Frame Hd/Ht Period        DP         n             L
     1    Ht    150  0.015200  1.587110   7627.929688
     1    Ht    200  0.007440  2.403570  94841.445312
     1    Ht    250  0.006170  2.309110 136895.203125
     1    Ht    300  0.006000  2.388160 147584.937500
```
