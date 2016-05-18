# Movielens ratings summary

Fast way to summarize ratings.csv from [movielens](http://files.grouplens.org/datasets/movielens/ml-latest-small-README.html). It is particular useful to process the [full](http://grouplens.org/datasets/movielens/latest/), contains 22M ratings (as of 05/16, the last update of this dataset was 1/2016).

Usage:

```sh
# Assuming the dataset was unzipped at /tmp
$ cd ${MY_WORKSPACE}
$ git clone https://github.com/danielfireman/phd/tree/master/movielens/ratings_summary
$ cd ratings_summary
$ cp header.csv /tmp/ml-latest/ratings_summ.csv && \
cat /tmp/ml-latest/ratings.csv | go run main.go  >> /tmp/ml-latest/ratings_summ.csv
```

Import those in R or your preferred tool and enjoy!