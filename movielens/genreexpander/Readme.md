# Movielens genre expander

Expands the [movielens](http://files.grouplens.org/datasets/movielens/ml-latest-small-README.html) genre field into columns.

Usage:

```sh
$ cp movies_header.csv movies_expanded.csv
$ cat movies.csv | go run main.go >> movies_expanded.csv
```

Import those in R or your preferred tool and enjoy!