# plfit-tablegen

A small utility to summarize `plfit` data.

## Install

```
git clone https://github.com/ivaivalous/plfit-tablegen.git
cd plfit-tablegen
go build
```

## Prerequisites
 - [`plfit`](https://pypi.org/project/plfit/) must be installed


## Run

```
plfit-tablegen <path containing .dat files> [--silent]
```

## Output

For example,

```
 Frame  Hd/Ht  Period         DP         n
     1      Ht     150  0.015200  1.587110
     1      Ht     200  0.007440  2.403570
     1      Ht     250  0.025170  2.909110
     1      Ht     300  0.006000  1.388160
     1      Ht     350  0.005700  2.941200
     1      Ht     400  0.025750  2.393290
     2      Ht     150  0.007650  2.523130
     2      Ht     200  0.042210  2.945640
     2      Ht     250  0.010400  1.337510
     2      Ht     300  0.004580  1.686440
     2      Ht     350  0.005180  1.653440
     2      Ht     400  0.013320  2.266300
```
